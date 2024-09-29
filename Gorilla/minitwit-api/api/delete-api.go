package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"minitwit-api/db"
	"minitwit-api/model"
	"minitwit-api/sim"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete handler invoked")
	db, err := db.GetDb()
	if err != nil {
		fmt.Println("Could not get database: ", err)
	}
	sim.UpdateLatest(r)
	dec := json.NewDecoder(r.Body)
	var rv model.DeleteData
	err = dec.Decode(&rv)
	if err != nil {
		fmt.Println("Error decoding request data: ", err)
	}
	is_auth := sim.Is_authenticated(w, r)
	if !is_auth {
		fmt.Println("Unauthorized access attempt to Delete")
		return
	}
	if rv.User != "" && r.Method == "POST" {
		toDeleteUsername := rv.User
		fmt.Println("Deleting user: ", toDeleteUsername)

		toDeleteUser_id, _ := db.Get_user_id(toDeleteUsername)
		fmt.Println("User ID to delete: ", toDeleteUser_id)

		db.QueryDelete([]int{toDeleteUser_id})
		fmt.Println("User deleted successfully")

	} else {
		fmt.Println("Invalid request: username missing or request method not POST")
	}

	fmt.Println("Delete completed")
	w.WriteHeader(http.StatusOK)
}
