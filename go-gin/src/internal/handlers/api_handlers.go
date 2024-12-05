package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go-gin/src/internal/db"
	"go-gin/src/internal/helpers"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// struct for error data (need to be JSON before return)
type ErrorData struct {
	status    int
	error_msg string
}

type UserData struct {
	Username string
	Email    string
	Pwd      string
}

type MessageData struct {
	Content string `json:"content"`
}

func Not_req_from_simulator(c *gin.Context) (statusCode int, errStr string) {
	auth := c.Request.Header.Get("Authorization")
	if auth != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		statusCode = 403
		errStr = "You are not authorized to use this resource!"
		return statusCode, errStr
	}
	return
}

func UpdateLatestHandler(c *gin.Context) {
	parsedCommandID := c.Query("latest")
	commandID, err := strconv.Atoi(parsedCommandID)

	if err != nil || parsedCommandID == "" {
		// Handle the case where the parameter is not present or cannot be converted to an integer
		commandID = -1
	}
	if commandID != -1 {
		err := db.UpdateLatest(commandID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update latest value"})
			return
		}
	}
}

func GetLatestHandler(c *gin.Context) {
	latestProcessedCommandID, err := db.GetLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read latest value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latest": latestProcessedCommandID})
}

func GetLatestHelper() int {
	latestProcessedCommandID, err := db.GetLatest()
	if err != nil {
		return -2
	}
	return latestProcessedCommandID
}

/*
/api/register
POST
Takes data from the POST and registers a user in the db
returns: ("", 204) or ({"status": 400, "error_msg": error}, 400)
*/
func ApiRegisterHandler(c *gin.Context) {
	UpdateLatestHandler(c)

	errorData := ErrorData{
		status:    0,
		error_msg: "",
	}

	// Check if user already exists
	userID, exists := c.Get("UserID")
	if exists {
		fmt.Println("Attempt to register an existing user")
		errorData.status = 400
		errorData.error_msg = "User already exists: " + fmt.Sprintf("%v", userID)
		c.AbortWithStatusJSON(400, errorData)
		return
	}

	if c.Request.Method == http.MethodPost {
		// Read the request body
		var registerReq UserData
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("Error failed to read request body")
			errorData.status = 400
			errorData.error_msg = "Failed to read JSON"
			c.AbortWithStatusJSON(400, errorData)
			return
		}

		// Parse the request body from JSON
		// Unmarshal parses the JSON and stores it in a pointer (registerReq)
		if err := json.Unmarshal(body, &registerReq); err != nil {
			fmt.Println("Failed to parse request JSON")
			errorData.status = 400
			errorData.error_msg = "Failed to parse JSON"
			c.AbortWithStatusJSON(400, errorData)
			return
		}

		// Set the user data
		username := registerReq.Username
		email := registerReq.Email
		password := registerReq.Pwd

		// Get user ID
		userID, err := db.GetUserIDByUsername(username)
		if err != nil {
			fmt.Println("Error getting username by id")
			errorData.status = 400
			errorData.error_msg = "Failed to get userID"
			c.AbortWithStatusJSON(400, errorData)
			return
		}

		// Check for errors
		if username == "" {
			errorData.status = 400
			errorData.error_msg = "You have to enter a username"
			c.AbortWithStatusJSON(400, errorData.error_msg)
			return

		} else if email == "" || !strings.Contains(email, "@") {
			errorData.status = 400
			errorData.error_msg = "You have to enter a valid email address"
			c.AbortWithStatusJSON(400, errorData.error_msg)
			return

		} else if password == "" {
			errorData.status = 400
			errorData.error_msg = "You have to enter a password"
			c.AbortWithStatusJSON(400, errorData.error_msg)
			return

		} else if fmt.Sprint(userID) != "-1" {
			errorData.status = 400
			errorData.error_msg = "The username is already taken"
			c.AbortWithStatusJSON(400, errorData.error_msg)
			return

		} else {
			hash := md5.Sum([]byte(password))
			err := db.RegisterUser(username, email, hash)
			if err != nil {
				fmt.Println("Failed registration attempt due to an error during registration")
				errorData.status = 400
				errorData.error_msg = "Failed to register user"
				c.AbortWithStatusJSON(400, errorData.error_msg)
				return
			}
		}

		if errorData.error_msg != "" {
			c.AbortWithStatusJSON(400, errorData.error_msg)
			return
		} else {
			c.JSON(204, "")
		}
	}
}

/*
/api/msgs
/api/msgs?no=<num>
*/
func ApiMsgsHandler(c *gin.Context) {
	UpdateLatestHandler(c)

	errorData := ErrorData{
		status:    0,
		error_msg: "",
	}

	not_req_from_sim_statusCode, not_req_from_sim_errStr := Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		errorData.status = http.StatusForbidden
		errorData.error_msg = not_req_from_sim_errStr
		c.AbortWithStatusJSON(http.StatusForbidden, errorData.error_msg)
		return
	}

	numMsgs := c.Request.Header.Get("no")
	numMsgsInt, err := strconv.Atoi(numMsgs)
	// fallback on default value
	if err != nil {
		fmt.Println("Falling back to default number of messages due to parsing error")
		numMsgsInt = 100
	}

	messages, err := db.GetPublicMessages(numMsgsInt)
	if err != nil {

		fmt.Println("Failed to fetch messages from DB")

		errorData.status = http.StatusBadRequest
		errorData.error_msg = "Failed to fetch messages from DB"
		c.AbortWithStatusJSON(http.StatusBadRequest, errorData)
	}

	filteredMessages := helpers.FilterMessages(messages)
	jsonFilteredMessages, _ := json.Marshal(filteredMessages)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(jsonFilteredMessages))
}

/*
/api/msgs/<username>
*/
func ApiMsgsPerUserHandler(c *gin.Context) {
	UpdateLatestHandler(c)

	errorData := ErrorData{
		status:    0,
		error_msg: "",
	}

	not_req_from_sim_statusCode, not_req_from_sim_errStr := Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		errorData.status = http.StatusForbidden
		errorData.error_msg = not_req_from_sim_errStr
		c.AbortWithStatusJSON(http.StatusForbidden, errorData.error_msg)
		return
	}

	profileUserName := c.Param("username")
	userId, err := db.GetUserIDByUsername(profileUserName)
	if userId == -1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if c.Request.Method == http.MethodGet {
		numMsgs := c.Request.Header.Get("no")
		numMsgsInt, err := strconv.Atoi(numMsgs)
		// fallback on default value
		if err != nil {
			numMsgsInt = 100
			fmt.Println("Fallback to default number of messages due to parsing error")
		}

		messages, err := db.GetUserMessages(userId, numMsgsInt)
		if err != nil {
			fmt.Println("Failed to fetch messages from DB")
			errorData.status = http.StatusBadRequest
			errorData.error_msg = "Failed to fetch messages from DB"
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorData)
		}

		// Log successful retrieval of messages
		filteredMessages := helpers.FilterMessages(messages)
		jsonFilteredMessages, _ := json.Marshal(filteredMessages)
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(jsonFilteredMessages))

	} else if c.Request.Method == http.MethodPost {
		// Read the request body
		var messageReq MessageData
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("Failed to read request body")
			errorData.status = 400
			errorData.error_msg = "Failed to read JSON"
			c.AbortWithStatusJSON(http.StatusBadRequest, errorData)
		}

		if err := json.Unmarshal(body, &messageReq); err != nil {
			fmt.Println("Failed to parse JSON body")
			errorData.status = 400
			errorData.error_msg = "Failed to parse JSON"
		}

		text := messageReq.Content
		authorId, err := db.GetUserIDByUsername(profileUserName)
		if err != nil {
			errorData.status = http.StatusBadRequest
			errorData.error_msg = "Failed to get userID"
			c.AbortWithStatusJSON(http.StatusBadRequest, errorData)
		}

		err = db.AddMessage(text, authorId)
		if err != nil {
			errorData.status = http.StatusInternalServerError
			errorData.error_msg = "Failed to upload message"
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorData)
		}

		c.String(http.StatusNoContent, "")
	}
}

/*
GET and POST
if GET:

	return: all followers that :username follows

else if POST:

	if FOLLOW:
		make userA follow userB
		return: status code
	if UNFOLLOW:
		make userA unfollow userB
		return: status code

/api/fllws/<username>
*/
func ApiFllwsHandler(c *gin.Context) {
	UpdateLatestHandler(c)

	errorData := ErrorData{
		status:    0,
		error_msg: "",
	}

	not_req_from_sim_statusCode, not_req_from_sim_errStr := Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		errorData.status = http.StatusForbidden
		errorData.error_msg = not_req_from_sim_errStr
		c.AbortWithStatusJSON(http.StatusForbidden, errorData.error_msg)
		return
	}

	if c.Request.Method == http.MethodGet {
		profileUserName := c.Param("username")
		numFollr := c.Request.Header.Get("no")
		numFollrInt, err := strconv.Atoi(numFollr)
		// fallback on default value
		if err != nil {
			fmt.Println("Fallback to default number of followers due to parsing error")
			numFollrInt = 100
		}

		userId, err := db.GetUserIDByUsername(profileUserName)
		if err != nil || userId == -1 {
			fmt.Println("Failed to get user ID for follow/unfollow actions")
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Fetch all followers for the user
		followers, err := db.GetFollowing(userId, numFollrInt)
		if err != nil {
			fmt.Println("Failed to fetch followers from DB")
			errorData.status = http.StatusInternalServerError
			errorData.error_msg = "Failed to fetch followers from DB"
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorData)
		}

		// empty slice for follower usernames
		followerNames := []string{}

		// Append the usernames to the followerNames slice
		for _, follower := range followers {
			followerNames = append(followerNames, string(follower.Username))
		}

		// Prepare response
		followersResponse := gin.H{
			"follows": followerNames,
		}

		// Send JSON response of all followers
		c.JSON(200, followersResponse)

	} else if c.Request.Method == http.MethodPost {
		// POST request
		var requestBody struct {
			Follow   string `json:"follow"`
			Unfollow string `json:"unfollow"`
		}

		// Bind JSON data to requestBody
		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println("Failed to bind JSON for follow/unfollow action")
			errorData.status = http.StatusNotFound
			errorData.error_msg = "Failed to parse JSON"
			c.AbortWithStatusJSON(http.StatusNotFound, errorData)
			return
		}

		profileUserName := c.Param("username")

		// Convert profileUserName to userID
		userId, err := db.GetUserIDByUsername(profileUserName)
		if err != nil || userId == -1 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if requestBody.Follow != "" {
			// Follow logic
			// Convert requestBody.Follow to profileUserID
			profileUserID, err := db.GetUserIDByUsername(requestBody.Follow)
			if err != nil || profileUserID == -1 {
				fmt.Println("Failed to get user ID for follow/unfollow actions")
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			// Follow the user
			if err := db.FollowUser(userId, profileUserID); err != nil {
				fmt.Println("Failed to follow user")
				errorData.status = http.StatusNotFound
				errorData.error_msg = "Failed to follow user"
				c.AbortWithStatusJSON(http.StatusNotFound, errorData)
				return
			}

			c.JSON(http.StatusNoContent, "")
			return
		} else if requestBody.Unfollow != "" {
			// Unfollow logic
			// Convert requestBody.Unfollow to profileUserID
			profileUserID, err := db.GetUserIDByUsername(requestBody.Unfollow)
			if err != nil || profileUserID == -1 {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			// Unfollow the user
			if err := db.UnfollowUser(userId, profileUserID); err != nil {
				fmt.Println("Failed to unfollow user")
				errorData.status = http.StatusNotFound
				errorData.error_msg = "Failed to unfollow user"
				c.AbortWithStatusJSON(http.StatusNotFound, errorData)
				return
			}

			c.JSON(http.StatusNoContent, "")
		} else {
			errorData.status = http.StatusNotFound
			errorData.error_msg = "No 'follow' or 'unfollow' provided in request"
			c.AbortWithStatusJSON(http.StatusNotFound, errorData)
			return
		}
	}
}
