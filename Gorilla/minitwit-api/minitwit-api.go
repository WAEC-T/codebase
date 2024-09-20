package main

import (
	"net/http"
	"os"

	"minitwit-api/api"
	"minitwit-api/db"
	"minitwit-api/db/postgres"
	"minitwit-api/logger"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var lg = logger.InitializeLogger()

func main() {
	lg.Info("Starting Minitwit API server")

	pgImpl := &postgres.PostgresDbImplementation{}
	pgImpl.Connect_db()
	db.SetDb(pgImpl)

	r := mux.NewRouter()

	r.HandleFunc("/health", api.Health).Name("Health")
	r.HandleFunc("/register", api.Register).Name("Register")
	r.HandleFunc("/msgs", api.Messages).Methods("GET").Name("Messages")
	r.HandleFunc("/msgs/{username}", api.Messages_per_user).Methods("GET", "POST").Name("Messages_per_user")
	r.HandleFunc("/fllws/{username}", api.Follow).Name("Follow")
	r.HandleFunc("/latest", api.Get_latest).Methods("GET").Name("Get_latest")
	r.HandleFunc("/cleandb", api.Cleandb).Name("Cleandb")
	r.HandleFunc("/delete", api.Delete).Name("Delete")

	r.Handle("/metrics", promhttp.Handler()).Name("Metrics")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "15001"
	}

	lg.Info("Listening on port:", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		lg.Fatal("Failed to start server: %v", err)
	}
}
