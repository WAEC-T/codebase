package main

import (
	"go-gin/src/internal/config"
	"go-gin/src/internal/db"
	"go-gin/src/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var err error

func main() {
	/*---------------------
	 * Load env vars
	 *----------------------*/
	uri := routes.LoadEnvVars()

	/*---------------------
	 * Connect to DB
	*----------------------*/
	config.DB, err = db.ConnectDB(uri) // Ensure this matches the function name in the db package
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	/*---------------------
	* Setup routing
	*----------------------*/
	r := gin.Default()
	routes.SetRouteHandlers(r)

}
