package handlers

import (
	"fmt"
	"go-gin/src/internal/db"
	"go-gin/src/internal/helpers"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		err := db.FollowUser(userID.(int), profileUserID)
		if err != nil {
			fmt.Println("follow user failed with:", err)
			c.Redirect(http.StatusFound, "/user/"+profileUserName)
			return
		}
		session.AddFlash("You are now following " + profileUserName)
	}
	if action == "/unfollow" {
		err := db.UnfollowUser(userID.(int), profileUserID)
		if err != nil {
			fmt.Println("unfollow user failed with:", err)
			c.Redirect(http.StatusFound, "/user/"+profileUserName)
			return
		}
		session.AddFlash("You are no longer following " + profileUserName)
	}

	if !helpers.SaveSessionOrRedirect(c, session.Save(), "/user/"+profileUserName) {
		return
	}
	c.Redirect(http.StatusFound, "/user/"+profileUserName)
}

func PublicTimelineHandler(c *gin.Context) {
	messages, err := db.GetPublicMessages(30)
	if err != nil {
		return
	}
	formattedMessages := helpers.FormatMessages(messages)

	context := gin.H{
		"TimelineBody": true,
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
	if userID == nil {
		c.Redirect(http.StatusFound, "/public") // Redirect to login page if session is invalid
		return
	}
	flashMessages := session.Flashes()

	if !helpers.SaveSessionOrRedirect(c, session.Save(), "/") {
		return
	}

	profileUserName := c.Param("username")
	profileUser, err := db.GetUserByUsername(profileUserName)

	if profileUser.Username == "" {
		fmt.Println("User not found for timeline")
		c.Redirect(http.StatusFound, "/public")
		return
	}
	if err != nil {
		if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
			fmt.Printf("Failed to abort with error: %v", errAbort)
		}
		return
	}

	// does the logged in user follow them
	pUserId := profileUser.UserID
	profileName := profileUser.Username
	userName, _ := db.GetUserNameByUserID(userID.(int))

	followed, _ := db.IsFollowing(userID.(int), pUserId)

	messages, err := db.GetUserMessages(pUserId, 30)

	if err != nil {
		if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
			fmt.Printf("Failed to abort with error: %v", errAbort)
		}
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
		c.Redirect(http.StatusFound, "/public") // Redirect to public timeline if session is invalid
		return
	}

	userName, err := db.GetUserNameByUserID(userID.(int))
	if err != nil {
		if err.Error() == "record not found" {
			c.Redirect(http.StatusSeeOther, "/public")
			return
		}
		if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
			fmt.Printf("Failed to abort with error: %v", errAbort)
		}
		return
	}

	flashMessages := session.Flashes()
	if !helpers.SaveSessionOrRedirect(c, session.Save(), "/") {
		return
	}

	messages, following, err := db.GetMyMessages(userID.(int))
	if err != nil {
		if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
			fmt.Printf("Failed to abort with error: %v", errAbort)
		}
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
		"Followed":     following,
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
			if !helpers.SaveSessionOrRedirect(c, session.Save(), "/") {
				return
			}
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
			if !helpers.SaveSessionOrRedirect(c, session.Save(), "/") {
				return
			}
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
		password2 := c.Request.FormValue("password2")

		userID, err := db.GetUserIDByUsername(userName)
		if err != nil {
			if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
				fmt.Printf("Failed to abort with error: %v", errAbort)
			}
			return
		}

		if userName == "" {
			errorData = "You have to enter a username"
		} else if email == "" || !strings.Contains(email, "@") {
			errorData = "You have to enter a valid email address"
		} else if password == "" {
			errorData = "You have to enter a password"
		} else if password != password2 {
			errorData = "The two passwords do not match"
		} else if fmt.Sprint(userID) != "-1" {
			errorData = "The username is already taken"
		} else {
			err := db.RegisterUser(userName, email, password)
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

			if !helpers.SaveSessionOrRedirect(c, session.Save(), "/register") {
				return
			}
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
	if !helpers.SaveSessionOrRedirect(c, session.Save(), "/login") {
		return
	}

	userID := session.Get("userID")
	if userID != nil {
		session.AddFlash("You were logged in")
		if !helpers.SaveSessionOrRedirect(c, session.Save(), "/login") {
			return
		}
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
			if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
				fmt.Printf("Failed to abort with error: %v", errAbort)
			}
			return
		}
		if user.Username == "" {
			errorData = "Invalid username"
		} else if !helpers.CheckPassword(password, user.Pwd) {
			errorData = "Invalid password"
		} else {
			// Save userID in the session
			session.Set("userID", user.UserID)
			session.AddFlash("You were logged in")
			if !helpers.SaveSessionOrRedirect(c, session.Save(), "/login") {
				return
			}

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
	if !helpers.SaveSessionOrRedirect(c, session.Save(), "/login") {
		return
	}

	// Redirect to the login page
	c.Redirect(http.StatusFound, "/login")
}
