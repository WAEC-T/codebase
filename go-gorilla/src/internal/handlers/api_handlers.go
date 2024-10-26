package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go-gorilla/src/internal/auth"
	"go-gorilla/src/internal/config"
	"go-gorilla/src/internal/db"
	"go-gorilla/src/internal/models"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func API_Follow(w http.ResponseWriter, r *http.Request) {
	//Get username
	vars := mux.Vars(r)
	username := vars["username"]

	//Update latest value
	API_UpdateLatestHandler(w, r)

	//Ensure authentication
	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Follow: ", username)
		return
	}

	//Get userID
	user_id, err := db.GetUserIDByUsername(username)
	if err != nil || user_id == -1 {
		fmt.Println("Error getting user ID", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//Convert user id to string
	userIDStr := strconv.Itoa(user_id)

	//Set follow-request type
	var rv models.FollowData

	// Decode JSON body for POST requests
	if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
			fmt.Println("Error decoding request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check if neither follow nor unfollow is provided
		if rv.Follow == "" && rv.Unfollow == "" {
			fmt.Println("Missing follow/unfollow field in request body")
			http.Error(w, "Missing 'follow' or 'unfollow' field", http.StatusBadRequest)
			return
		}

		// Check if it's a follow request
		if rv.Follow != "" {
			follow_username := rv.Follow
			follow_user_id, err := db.GetUserIDByUsername(follow_username)
			if err != nil || user_id == -1 {
				fmt.Println("Follow user not found or invalid user ID:", follow_username)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			// Convert follow_user_id to string and follow the user
			follower_userIDStr := strconv.Itoa(follow_user_id)
			db.FollowUser(userIDStr, follower_userIDStr)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Check if it's an unfollow request
		if rv.Unfollow != "" {
			unfollow_username := rv.Unfollow
			unfollow_user_id, err := db.GetUserIDByUsername(unfollow_username)
			if err != nil || user_id == -1 {
				fmt.Println("Unfollow user not found or invalid user ID:", unfollow_username)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			// Convert unfollow_user_id to string and unfollow the user
			unfollower_userIDStr := strconv.Itoa(unfollow_user_id)
			db.UnfollowUser(userIDStr, unfollower_userIDStr)
			w.WriteHeader(http.StatusNoContent)
			return
		}

	} else if r.Method == "GET" {
		followers, errx := db.GetFollowing(userIDStr, 100)
		if errx != nil {
			fmt.Println("Error getting followers for", username)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// empty slice for follower usernames
		followerNames := []string{}

		// Append the usernames to the followerNames slice
		for _, follower := range followers {
			followerNames = append(followerNames, follower.Username)
		}

		// Wrap the response in an object with a key "followers"
		response := map[string][]string{
			"follows": followerNames,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode the wrapped object
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println("Error encoding followers as JSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func API_GetLatestHandler(w http.ResponseWriter, r *http.Request) {
	count, err := db.GetLatest()
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Latest int `json:"latest"`
	}{
		Latest: count,
	})
}

func API_UpdateLatestHandler(w http.ResponseWriter, r *http.Request) {
	// Get the "latest" query parameter
	parsedCommandID := r.URL.Query().Get("latest")

	// Try to convert the parameter to an integer
	commandID, err := strconv.Atoi(parsedCommandID)

	if err != nil || parsedCommandID == "" {
		// Handle the case where the parameter is not present or cannot be converted to an integer
		commandID = -1
	}

	if commandID != -1 {
		// Attempt to update the latest command ID
		err := db.UpdateLatest(commandID)
		if err != nil {
			// Respond with a JSON error message
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update latest value"})
			return
		}
	}
}

func API_Messages(w http.ResponseWriter, r *http.Request) {
	//Update latest
	API_UpdateLatestHandler(w, r)

	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Messages")
		return
	}

	if r.Method == "GET" {
		messages, err := db.GetPublicMessages(100) //TODO: Agree on number

		if err != nil {
			fmt.Println("Error encoding JSON response:", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(messages)
	}
}

func API_Messages_per_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	//Update latest handler
	API_UpdateLatestHandler(w, r)

	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Messages_perUser")
		return
	}

	user_id, err := db.GetUserIDByUsername(username)
	if err != nil || user_id == -1 {
		fmt.Println("Error getting user ID", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		messages, err := db.GetUserMessages(user_id, 100)

		if err != nil {
			fmt.Println("Error encoding JSON response: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(messages)

	} else if r.Method == "POST" {
		var rv models.Messages

		err := json.NewDecoder(r.Body).Decode(&rv)
		fmt.Print(r)
		if err != nil {
			fmt.Println("Error decoding request body: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		db.AddMessage(rv.Content, user_id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func API_Register(w http.ResponseWriter, r *http.Request) {
	//Update latest
	API_UpdateLatestHandler(w, r)

	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Messages_perUser")
		return
	}

	var rv models.RegisterData
	err := json.NewDecoder(r.Body).Decode(&rv)
	if err != nil {
		fmt.Println("Error decoding request body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		hash := md5.Sum([]byte(rv.Pwd))
		db.RegisterUser(rv.Username, rv.Email, hash)
		w.WriteHeader(http.StatusNoContent)
	}
	if config.DB.Error != nil {
		fmt.Println("Registration error: ", config.DB.Error)
		w.WriteHeader(http.StatusBadRequest)
	}

}
