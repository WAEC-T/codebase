package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"gorilla-minitwit/internal/db"
	"gorilla-minitwit/internal/models"
	"gorilla-minitwit/internal/sim"
)

func Follow(w http.ResponseWriter, r *http.Request) {
	//Log
	fmt.Println("Follow handler invoked")

	//Get username
	vars := mux.Vars(r)
	username := vars["username"]

	//Update latest value
	sim.UpdateLatest(r)

	//Ensure authentication
	is_auth := sim.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Follow: ", username)
		return
	}

	//Get userID
	user_id, err := db.GetUserIDByUsername(username)
	if err != nil {
		fmt.Println("Error getting user ID", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//Convert user id to string
	userIDStr := strconv.Itoa(user_id)

	//Get number of followers -> TODO: how many?
	no_flws := No_followees(r)

	//Set follow-request type
	var rv models.FollowData

	//If the request is a POST but body is missing
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&rv)
		if err != nil {
			fmt.Println("Error decoding request body", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	//If the request is POST and it is a follow-request
	if r.Method == "POST" && rv.Follow != "" {
		follow_username := rv.Follow
		follow_user_id, err := db.GetUserIDByUsername(follow_username)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Follow user not found or invalid user ID", follow_username)
			return
		}

		// Perform type assertion for userID
		follower_userIDStr := strconv.Itoa(follow_user_id)

		db.FollowUser(userIDStr, follower_userIDStr)
		fmt.Println("User followed")
		w.WriteHeader(http.StatusNoContent)

		//If the request is POST and it is a unfollow-request
	} else if r.Method == "POST" && rv.Unfollow != "" {

		unfollow_username := rv.Unfollow
		unfollow_user_id, err := db.GetUserIDByUsername(unfollow_username)

		if err != nil {
			fmt.Println("Unfollow user not found or invalid user ID", unfollow_username)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		unfollower_userIDStr := strconv.Itoa(unfollow_user_id)

		db.UnfollowUser(userIDStr, unfollower_userIDStr)
		fmt.Println("User unfollowed")
		w.WriteHeader(http.StatusNoContent)

	} else if r.Method == "GET" {
		followers, errx := db.GetFollowing(userIDStr, no_flws)
		if errx != nil {
			fmt.Println("Error getting followers for", username)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Set response header as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and return followers as JSON
		if err := json.NewEncoder(w).Encode(followers); err != nil {
			fmt.Println("Error encoding followers as JSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("Retrieved followers for", username)
	}
}

func No_followees(r *http.Request) int {
	value := r.URL.Query().Get("no")
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return 100
}

func Get_latest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get latest handler invoked ")

	count := db.GetCount("sim")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Latest int `json:"latest"`
	}{
		Latest: count,
	})
}

func Messages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Messages handler invoked")
	db, err := db.GetDb()
	if err != nil {
		fmt.Println("Could not get database", err)
	}
	sim.UpdateLatest(r)

	is_auth := sim.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Messages")
		return
	}
	no_msg := no_msgs(r)

	if r.Method == "GET" {
		messages := db.GetMessages([]int{no_msg})

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(messages)

		if err != nil {
			fmt.Println("Error encoding JSON response:", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}

func Messages_per_user(w http.ResponseWriter, r *http.Request) {
	db, err := db.GetDb()
	if err != nil {
		fmt.Println("Could not get database - Messages per user", err)
	}
	vars := mux.Vars(r)
	username := vars["username"]
	sim.UpdateLatest(r)

	is_auth := sim.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Messages_perUser")
		return
	}
	no_msg := no_msgs(r)

	user_id, err := db.Get_user_id(username)
	if err != nil {
		fmt.Println("Error getting user ID", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method == "GET" {
		messages := db.GetMessagesForUser([]int{user_id, no_msg})

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			fmt.Println("Error encoding JSON response: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

	} else if r.Method == "POST" {
		var rv model.MessageData

		err := json.NewDecoder(r.Body).Decode(&rv)
		fmt.Print(r)
		if err != nil {
			fmt.Println("Error decoding request body: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		pubDate := strconv.FormatInt(time.Now().Unix(), 10)
		fmt.Println("pubdate: ", pubDate)
		message := &model.Messages{
			AuthorID: user_id,
			Text:     rv.Content,
			PubDate:  pubDate,
			Flagged:  0,
		}
		db.QueryMessage(message)
		fmt.Println("Message posted", user_id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func no_msgs(r *http.Request) int {
	value := r.URL.Query().Get("no")
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return 100
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register handler invoked")
	db, err := db.GetDb()
	if err != nil {
		fmt.Println("Could not get database: ", err)
	}
	sim.UpdateLatest(r)

	var rv model.RegisterData
	err = json.NewDecoder(r.Body).Decode(&rv)
	fmt.Println(r.Body)
	if err != nil {
		fmt.Println("Error decoding request body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		user_id, _ := db.Get_user_id(rv.Username)

		errMsg := ""

		if rv.Username == "" {
			errMsg = "You have to enter a username"
		} else if rv.Email == "" || !strings.Contains(rv.Email, "@") {
			errMsg = "You have to enter a valid email address"
		} else if rv.Pwd == "" {
			errMsg = "You have to enter a password"
		} else if !db.IsZero(user_id) {
			errMsg = "The username is already taken"
		} else {
			hash_pw := hashPassword(rv.Pwd)
			db.QueryRegister([]string{rv.Username, rv.Email, hash_pw})
			fmt.Println("User registered successfully", rv.Username)
			w.WriteHeader(http.StatusNoContent)
		}
		if errMsg != "" {
			fmt.Println("Registration error: ", errMsg)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func hashPassword(password string) string {
	return fmt.Sprintf("%d", xxhash.Sum64([]byte(password)))
}
