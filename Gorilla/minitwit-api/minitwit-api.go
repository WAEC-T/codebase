package main

import (
	"fmt"
	"minitwit-api/api"
	"minitwit-api/db"
	"minitwit-api/db/postgres"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("Starting Minitwit API server")

	pgImpl := &postgres.PostgresDbImplementation{}
	pgImpl.Connect_db()
	db.SetDb(pgImpl)

	r := mux.NewRouter()

	r.HandleFunc("/register", api.Register).Name("Register")
	r.HandleFunc("/msgs", api.Messages).Methods("GET").Name("Messages")
	r.HandleFunc("/msgs/{username}", api.Messages_per_user).Methods("GET", "POST").Name("Messages_per_user")
	r.HandleFunc("/fllws/{username}", api.Follow).Name("Follow")
	r.HandleFunc("/latest", api.Get_latest).Methods("GET").Name("Get_latest")
	//r.HandleFunc("/cleandb", api.Cleandb).Name("Cleandb")
	//r.HandleFunc("/delete", api.Delete).Name("Delete")

	r.Handle("/metrics", promhttp.Handler()).Name("Metrics")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "15001"
	}

	fmt.Println("Listening on port:", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println("Failed to start server: %v", err)
	}
}
