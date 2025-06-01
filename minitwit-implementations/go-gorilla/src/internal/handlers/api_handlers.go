package handlers

import (
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
	//Update latest value
	API_UpdateLatestHandler(w, r)

	//Ensure authentication
	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Request denied: not from simulator")
		return
	}

	//Get username
	vars := mux.Vars(r)
	username := vars["username"]

	//Get userID
	user_id, err := db.GetUserIDByUsername(username)
	if err != nil || user_id == -1 {
		fmt.Println("Error getting user ID", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//Set follow-request type
	var rv models.FollowData

	if r.Method == "POST" {

		// Decode JSON body for POST requests
		if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
			fmt.Println("Error decoding request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check if it's a follow request
		if rv.Follow != "" {
			follow_user_id, err := db.GetUserIDByUsername(rv.Follow)
			if err != nil || user_id == -1 {
				fmt.Println("Failed to get user ID for follow/unfollow actions")
				w.WriteHeader(http.StatusNotFound)
				return
			}

			// Follow the user
			if err := db.FollowUser(user_id, follow_user_id); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Check if it's an unfollow request
		if rv.Unfollow != "" {
			unfollow_user_id, err := db.GetUserIDByUsername(rv.Unfollow)
			if err != nil || user_id == -1 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			// Unfollow the user
			if err := db.UnfollowUser(user_id, unfollow_user_id); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

	} else if r.Method == "GET" {
		followers, errx := db.GetFollowing(user_id, 100)
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
		messages, err := db.GetPublicMessages(100)

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
	//Update latest handler
	API_UpdateLatestHandler(w, r)

	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Messages_perUser")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

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
		if err != nil {
			fmt.Println("Error decoding request body: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err = db.AddMessage(rv.Content, user_id)
		if err != nil {
			fmt.Println("Error getting user ID", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func API_Register(w http.ResponseWriter, r *http.Request) {
	//Update latest
	API_UpdateLatestHandler(w, r)

	is_auth := auth.Is_authenticated(w, r)
	if !is_auth {
		w.WriteHeader(http.StatusBadRequest)
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

	user, err := db.GetUserByUsername(rv.Username)
	if err == nil && user.Username != "" {
		fmt.Println("User already exists: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		db.RegisterUser(rv.Username, rv.Email, rv.Pwd)
		w.WriteHeader(http.StatusNoContent)
	}
	if config.DB.Error != nil {
		fmt.Println("Registration error: ", config.DB.Error)
		w.WriteHeader(http.StatusBadRequest)
	}

}
