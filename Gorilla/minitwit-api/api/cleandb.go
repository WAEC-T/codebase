package api

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"minitwit-api/db"
	"minitwit-api/sim"
)

func Cleandb(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Cleandb handler invoked")
	db, err := db.GetDb()
	if err != nil {
		fmt.Println("Could not get database: ", err)
	}
	is_auth := sim.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Cleandb")
		return
	}
	user_ids := make([]int, 4)
	usernames := []string{"a", "b", "c", "test"}

	for i, username := range usernames {
		user_id, _ := db.Get_user_id(username)
		user_ids[i] = user_id
		fmt.Println("Retrieved user ID for username: ", username, ", ID: ", user_id)
	}

	for _, userID := range user_ids {
		if !db.IsZero(userID) {
			db.QueryDelete([]int{userID})
			fmt.Println("Deleted user with ID: ", userID)
		} else {
			fmt.Println("Skipping deletion for user ID: ", userID)
		}
	}
	fmt.Println("Cleandb completed successfully")
	w.WriteHeader(http.StatusOK)
}
