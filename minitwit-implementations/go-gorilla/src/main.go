package main

import (
	"go-gorilla/src/internal/config"
	"go-gorilla/src/internal/db"
	"go-gorilla/src/internal/routes"
	"log"
	"net/http"
	"text/template"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var err error

func main() {
	/*----------------------
	 * Load env vars
	 *----------------------*/
	uri := routes.LoadEnvVars()

	/*-----------------------
	 * Connect to DB
	 *----------------------*/
	config.DB, err = db.ConnectDB(uri) // Ensure this matches the function name in the db package
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	/*---------------------
	* Setup routing
	*----------------------*/
	funcMap := routes.SetupRouting()

	config.Tpl, err = template.New("timeline.html").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	/*---------------------
	* Setup route-handlers
	*----------------------*/
	r := mux.NewRouter()
	routes.SetRouteHandlers(r)
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
