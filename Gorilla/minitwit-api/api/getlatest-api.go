package api

import (
	"encoding/json"
	"fmt"
	"minitwit-api/db"
	"net/http"
)

func Get_latest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get latest handler invoked ")
	db, err := db.GetDb()
	if err != nil {
		fmt.Println("Could not get database", err)
	}
	count := db.GetCount("sim")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Latest int `json:"latest"`
	}{
		Latest: count,
	})
}
