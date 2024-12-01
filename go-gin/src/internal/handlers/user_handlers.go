package handlers

import (
	"crypto/md5"
	"fmt"
	"go-gin/src/internal/db"
	"go-gin/src/internal/helpers"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Handlers
func UserFollowActionHandler(c *gin.Context) {
	session := sessions.Default(c)

	userID := session.Get("userID")

	profileUserName := c.Param("username")
	profileUser, err := db.GetUserByUsername(profileUserName)
	if err != nil {
		fmt.Println("get user failed with:", err)
		c.Redirect(http.StatusFound, "/public")
		return
	}
	profileUserID := profileUser.UserID

	action := c.Param("action")

	if action == "/follow" {
		db.FollowUser(userID.(int), profileUserID)
		session.AddFlash("You are now following " + profileUserName)
	}
	if action == "/unfollow" {
		db.UnfollowUser(userID.(int), profileUserID)
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

	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID != nil {
		userName, errName := db.GetUserNameByUserID(userID.(int))
		if errName == nil {
			context["UserName"] = userName
			context["UserID"] = userID.(int)
		}
	}

	// Render timeline template with the context including link variables
	c.HTML(http.StatusOK, "timeline.html", context)
}

func UserTimelineHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
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
	followed, _ := db.GetFollowing(userID.(int), 30) //TODO: LIMIT OF FOLLOWERS WE QUERY?
	pUserId := profileUser.UserID
	profileName := profileUser.Username
	userName, _ := db.GetUserNameByUserID(userID.(int))

	messages, err := db.GetUserMessages(pUserId, 30)

	if err != nil {
		fmt.Println("Error fetching user messages")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	formattedMessages := helpers.FormatMessages(messages)

	c.HTML(http.StatusOK, "timeline.html", gin.H{
		"TimelineBody":    true,
		"Endpoint":        "user_timeline",
		"UserID":          userID.(int),
		"UserName":        userName,
		"Messages":        formattedMessages,
		"Followed":        followed,
		"ProfileUser":     pUserId,
		"ProfileUserName": profileName,
		"Flashes":         flashMessages,
	})
}

func MyTimelineHandler(c *gin.Context) {
	session := sessions.Default(c)

	userID := session.Get("userID")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login") // Redirect to login page if session is invalid
		return
	}

	userName, err := db.GetUserNameByUserID(userID.(int))
	if err != nil {
		fmt.Println("Error getting username by id")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	flashMessages := session.Flashes()
	session.Save() // Clear flashes after retrieving

	messages, err := db.GetMyMessages(userID.(int))
	if err != nil {
		fmt.Println("Error getting users messages")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	formattedMessages := helpers.FormatMessages(messages)

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
	})
}

func AddMessageHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")

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
			err := db.AddMessage(text, userID.(int))
			if err != nil {
				fmt.Println("Error failed to add message")
				errorData = "Failed to add message"
				c.Redirect(http.StatusInternalServerError, "/?error="+errorData)
				return
			}

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

			// Redirect to login page after successful registration
			session.AddFlash("You were successfully registered and can login now")

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

	userID := session.Get("userID")
	if userID != nil {
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

			// Save userID in the session
			session.Set("userID", userID)
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

	// Clear all session data
	session.Clear()
	// Optionally add a flash message before clearing the session
	session.AddFlash("You were logged out")

	// Save the session to apply changes
	session.Save()

	// Redirect to the login page
	c.Redirect(http.StatusFound, "/login")
}
