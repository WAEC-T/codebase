package api

import (
	"encoding/json"
	"fmt"
	"minitwit-api/db"
	"minitwit-api/model"
	"minitwit-api/sim"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

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
