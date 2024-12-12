package main

import (
	"go-gin/src/internal/config"
	"go-gin/src/internal/db"
	"go-gin/src/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

var err error

func main() {

	//Set to ReleaseMode to disable logging
	gin.SetMode(gin.ReleaseMode)
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
	r := gin.New()
	routes.SetRouteHandlers(r)

	/*---------------------
	* Start the server
	*----------------------*/
	err = r.Run(":5000") // Start the server on port 5000
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
