package handlers

import (
	"encoding/json"
	"fmt"
	"go-gin/src/internal/auth"
	"go-gin/src/internal/config"
	"go-gin/src/internal/db"
	"go-gin/src/internal/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserData struct {
	Username string
	Email    string
	Pwd      string
}

type MessageData struct {
	Content string `json:"content"`
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
	//Update latest
	UpdateLatestHandler(c)

	not_req_from_sim_statusCode, not_req_from_sim_errStr := auth.Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		c.AbortWithStatusJSON(http.StatusForbidden, "Request denied: not from simulator")
		return
	}

	// Check if user already exists
	//TODO: DO IT

	// Read the request body
	var registerReq UserData
	err := json.NewDecoder(c.Request.Body).Decode(&registerReq)
	if err != nil {
		fmt.Println("Error failed to read request body")
		c.AbortWithStatusJSON(400, "Error failed to read request body")
		return
	}

	// Get user ID
	user, err := db.GetUserByUsername(registerReq.Username)
	if err == nil && user.Username != "" {
		fmt.Println("Error getting username by id")
		c.AbortWithStatusJSON(400, "Error getting username by id")
		return
	}

	if c.Request.Method == http.MethodPost {
		err := db.RegisterUser(registerReq.Username, registerReq.Email, registerReq.Pwd)
		if err != nil {
			fmt.Println("Failed registration attempt due to an error during registration")
			c.AbortWithStatusJSON(400, "Failed to register user")
			return
		}
		if config.DB.Error != nil {
			fmt.Println("Registration error: ", config.DB.Error)
			c.AbortWithStatusJSON(http.StatusBadRequest, "Registration error: ")
		}
	}

	c.JSON(204, "")
}

/*
/api/msgs
/api/msgs?no=<num>
*/
func ApiMsgsHandler(c *gin.Context) {
	UpdateLatestHandler(c)

	not_req_from_sim_statusCode, not_req_from_sim_errStr := auth.Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		c.AbortWithStatusJSON(http.StatusForbidden, "Request denied: not from simulator")
		return
	}

	if c.Request.Method == http.MethodGet {
		messages, err := db.GetPublicMessages(100)
		if err != nil {
			fmt.Println("Failed to fetch messages from DB")
			c.AbortWithStatusJSON(http.StatusBadRequest, "Failed to fetch messages from DB")
		}

		filteredMessages := helpers.FilterMessages(messages)
		jsonFilteredMessages, _ := json.Marshal(filteredMessages)
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(jsonFilteredMessages))
	}
}

/*
/api/msgs/<username>
*/
func ApiMsgsPerUserHandler(c *gin.Context) {
	UpdateLatestHandler(c)

	not_req_from_sim_statusCode, not_req_from_sim_errStr := auth.Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		c.AbortWithStatusJSON(http.StatusForbidden, "Request denied: not from simulator")
		return
	}

	profileUserName := c.Param("username")
	userId, err := db.GetUserIDByUsername(profileUserName)
	if userId == -1 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil || userId == -1 {
		if errAbort := c.AbortWithError(http.StatusInternalServerError, err); errAbort != nil {
			fmt.Printf("Failed to abort with error: %v", errAbort)
		}
		return
	}

	if c.Request.Method == http.MethodGet {
		messages, err := db.GetUserMessages(userId, 100)
		if err != nil {
			fmt.Println("Failed to fetch messages from DB")
			c.AbortWithStatusJSON(http.StatusInternalServerError, "Failed to fetch messages from DB")
		}

		filteredMessages := helpers.FilterMessages(messages)
		jsonFilteredMessages, _ := json.Marshal(filteredMessages)
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(jsonFilteredMessages))

	} else if c.Request.Method == http.MethodPost {
		var messageReq MessageData

		err := json.NewDecoder(c.Request.Body).Decode(&messageReq)
		if err != nil {
			fmt.Println("Failed to read request body")
			c.AbortWithStatusJSON(http.StatusBadRequest, "Failed to read request body")
			return
		}

		err = db.AddMessage(messageReq.Content, userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "Failed to upload message")
		}

		c.JSON(http.StatusNoContent, "")
	}
}

func ApiFllwsHandler(c *gin.Context) {
	//Update latest value
	UpdateLatestHandler(c)

	//Ensure authentication
	not_req_from_sim_statusCode, not_req_from_sim_errStr := auth.Not_req_from_simulator(c)
	if not_req_from_sim_statusCode == 403 && not_req_from_sim_errStr != "" {
		fmt.Println("Request denied: not from simulator")
		return
	}

	//Get username
	profileUserName := c.Param("username")

	//Get userID
	userId, err := db.GetUserIDByUsername(profileUserName)
	if err != nil || userId == -1 {
		fmt.Println("Failed to get user ID for follow/unfollow actions")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//Set follow-request type
	var requestBody struct {
		Follow   string `json:"follow"`
		Unfollow string `json:"unfollow"`
	}

	if c.Request.Method == http.MethodPost {

		// Decode JSON body for POST requests
		if err := json.NewDecoder(c.Request.Body).Decode(&requestBody); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, "Failed decoding request")
			return
		}

		// Check if neither follow nor unfollow is provided
		if requestBody.Follow == "" && requestBody.Unfollow == "" {
			fmt.Println("Missing follow/unfollow field in request body")
			c.AbortWithStatusJSON(http.StatusNotFound, "Missing follow/unfollow field in request body")
			return
		}

		if requestBody.Follow != "" {
			follow_user_id, err := db.GetUserIDByUsername(requestBody.Follow)
			if err != nil || follow_user_id == -1 {
				fmt.Println("Failed to get user ID for follow/unfollow actions")
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			// Follow the user
			if err := db.FollowUser(userId, follow_user_id); err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, "Failed to follow user")
				return
			}

			c.JSON(http.StatusNoContent, "")
			return

		} else if requestBody.Unfollow != "" {
			unfollow_user_id, err := db.GetUserIDByUsername(requestBody.Unfollow)
			if err != nil || unfollow_user_id == -1 {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			// Unfollow the user
			if err := db.UnfollowUser(userId, unfollow_user_id); err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, "Failed to unfollow user")
				return
			}

			c.JSON(http.StatusNoContent, "")
		}

	} else if c.Request.Method == http.MethodGet {
		followers, err := db.GetFollowing(userId, 100)
		if err != nil {
			fmt.Println("Failed to fetch followers from DB")
			c.AbortWithStatusJSON(http.StatusInternalServerError, "Failed to fetch followers from DB")
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

		c.JSON(200, followersResponse)
	}
}
