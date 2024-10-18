package routes

import (
	"go-minitwit/src/internal/handlers"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetRouteHandlers(r *gin.Engine) {

	r.LoadHTMLGlob("./templates/*.html")

	// sessions, for cookies
	store := cookie.NewStore([]byte("devops"))
	r.Use(sessions.Sessions("session", store))

	// Static (styling)
	r.Static("/static", "./static")

	// Define routes -> Here is where the links are being registered! Check the html layout file
	// user routes
	r.GET("/", handlers.MyTimelineHandler)
	r.GET("/public", handlers.PublicTimelineHandler)
	r.GET("/:username", handlers.UserTimelineHandler)
	r.GET("/register", handlers.RegisterHandler)
	r.GET("/login", handlers.LoginHandler)
	r.GET("/logout", handlers.LogoutHandler)
	r.GET("/:username/*action", handlers.UserFollowActionHandler)

	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/add_message", handlers.AddMessageHandler)

	// // API routes
	// // is it easier to separate the next two routes into two handlers?
	// r.GET("/api/msgs", apiMsgsHandler)
	// r.GET("/api/msgs/:username", apiMsgsPerUserHandler)
	// r.GET("/api/fllws/:username", apiFllwsHandler)

	// r.POST("/api/register", apiRegisterHandler)
	// r.POST("/api/msgs/:username", apiMsgsPerUserHandler)
	// r.POST("/api/fllws/:username", apiFllwsHandler)

	// // some helper method to "cache" what was the latest simulator action
	// r.GET("/api/latest", getLatestHandler)

	// Start the server
	r.Run(":8081")
}

func LoadEnvVars() string {
	// Save env-vars
	dbURL := os.Getenv("DATABASE_URL")

	return dbURL
}
