package handlers

import (
	"crypto/md5"
	"fmt"
	"gin-minitwit/src/internal/db"
	"gin-minitwit/src/internal/helpers"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Handlers
func UserFollowActionHandler(c *gin.Context) {
	session := sessions.Default(c)

	userID, errID := c.Cookie("UserID")
	if errID != nil {
		session.AddFlash("You need to login before continuing to follow or unfollow.")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return

	}
	profileUserName := c.Param("username")
	profileUser, err := db.GetUserByUsername(profileUserName)
	if err != nil {
		fmt.Println("get user failed with:", err)
		c.Redirect(http.StatusFound, "/public")
		return
	}
	profileUserID := fmt.Sprintf("%v", profileUser.UserID)

	action := c.Param("action")

	if action == "/follow" {
		db.FollowUser(userID, profileUserID)
		session.AddFlash("You are now following " + profileUserName)
	}
	if action == "/unfollow" {
		db.UnfollowUser(userID, profileUserID)
		session.AddFlash("You are no longer following " + profileUserName)
	}
	session.Save()
	c.Redirect(http.StatusFound, "/"+profileUserName)
}

func PublicTimelineHandler(c *gin.Context) {
	// need to pass a default value to getPublicMessages (GoLang doesn't support default values for arguments)
	messages, err := db.GetPublicMessages(30) //30 per page
	if err != nil {
		return
	}
	formattedMessages := helpers.FormatMessages(messages)

	context := gin.H{
		"TimelineBody": true, // This seems to be a flag you use to render specific parts of your layout
		"Endpoint":     "public_timeline",
		"Messages":     formattedMessages,
	}

	userID, errID := c.Cookie("UserID")
	if errID == nil {
		context["UserID"] = userID
		userName, errName := db.GetUserNameByUserID(userID)

		if errName == nil {
			context["UserName"] = userName
		}
	}

	// Render timeline template with the context including link variables
	c.HTML(http.StatusOK, "timeline.html", context)
}

func UserTimelineHandler(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()
	profileUserName := c.Param("username")
	profileUser, err := db.GetUserByUsername(profileUserName)

	if profileUser.Username == "" {
		fmt.Println("User not found for timeline")
		c.AbortWithStatus(404)
		return
	}
	if err != nil {
		fmt.Println("Error fetching user for timeline")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// does the logged in user follow them
	followed := false
	pUserId := profileUser.UserID
	profileName := profileUser.Username
	userID, errID := c.Cookie("UserID")
	userIDInt, _ := strconv.Atoi(userID)
	userName, _ := db.GetUserNameByUserID(userID)

	if errID == nil {
		followed, err = db.CheckFollowStatus(userIDInt, pUserId)
		if err != nil {
			fmt.Println("Error checking follow status")
			return
		}
	}

	messages, err := db.GetUserMessages(pUserId, 30)
	fmt.Println(messages)

	if err != nil {
		fmt.Println("Error fetching user messages")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	formattedMessages := helpers.FormatMessages(messages)
	fmt.Println("Rendering users public timeline")

	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"TimelineBody":    true,
		"Endpoint":        "user_timeline",
		"UserID":          userIDInt,
		"UserName":        userName,
		"Messages":        formattedMessages,
		"Followed":        followed,
		"ProfileUser":     pUserId,
		"ProfileUserName": profileName,
		"Flashes":         flashMessages,
	})
}

func MyTimelineHandler(c *gin.Context) {
	userID, err := c.Cookie("UserID")
	errMsg := c.Query("error")

	if err != nil {
		fmt.Println("Error getting user information")
		c.Redirect(http.StatusFound, "/public")
		return
	}

	userName, err := db.GetUserNameByUserID(userID)
	if err != nil {
		fmt.Println("Error getting username by id")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save() // Clear flashes after retrieving

	messages, err := db.GetMyMessages(userID)
	if err != nil {
		fmt.Println("Error getting users messages")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	formattedMessages := helpers.FormatMessages(messages)
	fmt.Println(formattedMessages)

	fmt.Println("Rendering users timeline")

	// For template rendering with Gin
	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"TimelineBody": true,
		"Endpoint":     "my_timeline",
		"UserID":       userID,
		"UserName":     userName,
		"Messages":     formattedMessages,
		"Followed":     false,
		"ProfileUser":  userID,
		"Flashes":      flashMessages,
		"Error":        errMsg,
	})
}

func AddMessageHandler(c *gin.Context) {
	session := sessions.Default(c)

	userID, err := c.Cookie("UserID")
	userIDString, errStr := strconv.Atoi(userID)
	if err != nil || errStr != nil {
		fmt.Println("Error getting user information")
		c.Redirect(http.StatusFound, "/public")
		return
	}

	var errorData string
	if c.Request.Method == http.MethodPost {
		err := c.Request.ParseForm()
		if err != nil {
			fmt.Println("Error failed to parse form data")
			errorData = "Failed to parse form data"
			c.Redirect(http.StatusBadRequest, "/?error="+errorData)
			return
		}

		// Validate form data
		text := c.Request.FormValue("text")

		if text == "" {
			c.Redirect(http.StatusSeeOther, "/")
			session.AddFlash("You have to enter a value")
			session.Save()
			return
		} else {
			err := db.AddMessage(text, userIDString)
			if err != nil {
				fmt.Println("Error failed to add message")
				errorData = "Failed to add message"
				c.Redirect(http.StatusInternalServerError, "/?error="+errorData)
				return
			}

			fmt.Println("Rendering users timeline")

			c.Redirect(http.StatusSeeOther, "/")
			session.AddFlash("Your message was recorded")
			session.Save()
			return
		}
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func RegisterHandler(c *gin.Context) {
	session := sessions.Default(c)

	userID, exists := c.Get("UserID")
	if exists {
		fmt.Println("User exists")
		fmt.Println("userID:", userID)
		return
	}

	var errorData string
	if c.Request.Method == http.MethodPost {
		err := c.Request.ParseForm()
		if err != nil {
			fmt.Println("Error failed to parse form data")
			errorData = "Failed to parse form data"
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"RegisterBody": true,
				"Error":        errorData,
			})
			return
		}

		// Validate form data
		userName := c.Request.FormValue("username")
		email := c.Request.FormValue("email")
		password := c.Request.FormValue("password")
		passwordConfirm := c.Request.FormValue("passwordConfirm")

		userID, err := db.GetUserIDByUsername(userName)
		if err != nil {
			fmt.Println("Error getting username by id")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if userName == "" {
			errorData = "You have to enter a username"
		} else if email == "" || !strings.Contains(email, "@") {
			errorData = "You have to enter a valid email address"
		} else if password == "" {
			errorData = "You have to enter a password"
		} else if password != passwordConfirm {
			errorData = "The two passwords do not match"
		} else if fmt.Sprint(userID) != "-1" {
			errorData = "The username is already taken"
		} else {
			hash := md5.Sum([]byte(password))
			err := db.RegisterUser(userName, email, hash)
			if err != nil {
				fmt.Println("Failed registration attempt due to an error during registration")
				errorData = "Failed to register user"
				c.HTML(http.StatusInternalServerError, "register.html", gin.H{
					"RegisterBody": true,
					"Error":        errorData,
				})
				return
			}

			fmt.Println("User successfully registered")

			// Redirect to login page after successful registration
			session.AddFlash("You were successfully registered and can login now")
			// print session info
			fmt.Println("session info:", session, "Logged in")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
	}
	c.HTML(http.StatusOK, "register.html", gin.H{
		"RegisterBody": true,
		"Error":        errorData,
	})
}

func LoginHandler(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()

	userID, _ := c.Cookie("UserID")
	if userID != "" {
		fmt.Println("User already logged in, redirecting")
		session.AddFlash("You were logged in")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	var errorData string

	if c.Request.Method == http.MethodPost {
		err := c.Request.ParseForm()
		if err != nil {
			fmt.Println("Failed to parse login form")
			errorData = "Failed to parse form data"
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"loginBody": true,
				"Error":     errorData,
			})
			return
		}

		userName := c.Request.FormValue("username")
		password := c.Request.FormValue("password")

		user, err := db.GetUserByUsername(userName)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if user.Username == "" {
			errorData = "Invalid username"
		} else if !helpers.CheckPasswordHash(password, user.PwHash) {
			errorData = "Invalid password"
		} else {
			userID, err := db.GetUserIDByUsername(userName)

			if err != nil {
				fmt.Println("Failed to retrieve userID during login")
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			fmt.Println("User successfully logged in")
			c.SetCookie("UserID", fmt.Sprint(userID), 3600, "/", "", false, true)

			session.AddFlash("You were logged in")
			session.Save()
			c.Redirect(http.StatusFound, "/")
			return
		}

	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"LoginBody": true,
		"Error":     errorData,
		"Flashes":   flashMessages,
	})
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println("User successfully logged out")

	session.AddFlash("You were logged out")

	session.Save()
	// Invalidate the cookie by setting its max age to -1
	// will delete the cookie <- nice stuff
	c.SetCookie("UserID", "", -1, "/", "", false, true)
	// redirect the user to the home page or login page
	c.Redirect(http.StatusFound, "/login")
}
