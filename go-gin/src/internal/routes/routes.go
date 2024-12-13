package routes

import (
	"go-gin/src/internal/handlers"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetRouteHandlers(r *gin.Engine) {

	r.LoadHTMLGlob("templates/*.html")

	// sessions, for cookies
	store := cookie.NewStore([]byte("devops"))
	r.Use(sessions.Sessions("session", store))

	// Static (styling)
	r.Static("static", "static")

	// Define routes -> Here is where the links are being registered! Check the html layout file
	// user routes
	r.GET("/", handlers.MyTimelineHandler)
	r.GET("/public", handlers.PublicTimelineHandler)
	r.GET("/user/:username", handlers.UserTimelineHandler)
	r.GET("/register", handlers.RegisterHandler)
	r.GET("/login", handlers.LoginHandler)
	r.GET("/logout", handlers.LogoutHandler)
	r.GET("/:username/*action", handlers.UserFollowActionHandler)

	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/add_message", handlers.AddMessageHandler)

	// API routes
	r.GET("/api/msgs", handlers.ApiMsgsHandler)
	r.GET("/api/msgs/:username", handlers.ApiMsgsPerUserHandler)
	r.GET("/api/fllws/:username", handlers.ApiFllwsHandler)

	r.POST("/api/register", handlers.ApiRegisterHandler)
	r.POST("/api/msgs/:username", handlers.ApiMsgsPerUserHandler)
	r.POST("/api/fllws/:username", handlers.ApiFllwsHandler)

	// some helper method to "cache" what was the latest simulator action
	r.GET("/api/latest", handlers.GetLatestHandler)
}

func LoadEnvVars() string {
	// Save env-vars
	dbURL := os.Getenv("DATABASE_URL")

	return dbURL
}
