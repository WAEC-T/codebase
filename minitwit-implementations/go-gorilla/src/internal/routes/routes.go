package routes

import (
	"go-gorilla/src/internal/handlers"
	"go-gorilla/src/internal/helpers"
	"go-gorilla/src/internal/models"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

func SetupRouting() template.FuncMap {
	funcMap := template.FuncMap{
		"getavatar": func(url string, size int) string {
			return handlers.Gravatar_url(url, size)
		},
		"gettimestamp": func(timestamp string) string {
			return helpers.Format_datetime(timestamp)
		},
		"url_for": func(routename, username string) string {
			switch routename {
			case "unfollow":
				return "/" + username + "/unfollow"
			case "follow":
				return "/" + username + "/follow"
			case "add_message":
				return "/add_message"
			case "my_timeline":
				return "/" 
			case "public_timeline":
				return "/public"
			case "user_timeline":
				return "/user/" + username
			case "logout":
				return "/logout"
			case "login":
				return "/login"
			case "register":
				return "/register"
			default:
				return "/"
			}
		},
		"formatUsernameUrl": func(username string) string {
			return strings.Replace(username, " ", "%20", -1)
		},
		"IsFollowing": func(following []models.Users, messageAuthorId int) bool {
			for _, user := range following {
				if user.UserID == messageAuthorId {
					return true
				}
			}
			return false
		},
	}
	return funcMap
}

func SetRouteHandlers(r *mux.Router) {
	//UI
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/public", handlers.Public_timeline)
	r.HandleFunc("/register", handlers.Register)
	r.HandleFunc("/login", handlers.Login)
	r.HandleFunc("/logout", handlers.Logout)
	r.HandleFunc("/", handlers.MyTimeline)
	r.HandleFunc("/add_message", handlers.Add_message).Methods("POST")
	r.HandleFunc("/{username}/follow", handlers.Follow_user)
	r.HandleFunc("/user/{username}", handlers.User_timeline)
	r.HandleFunc("/{username}/unfollow", handlers.Unfollow_user)

	//API
	r.HandleFunc("/api/msgs", handlers.API_Messages).Methods("GET").Name("Messages")
	r.HandleFunc("/api/msgs/{username}", handlers.API_Messages_per_user).Methods("GET", "POST").Name("Messages per user")
	r.HandleFunc("/api/fllws/{username}", handlers.API_Follow).Methods("GET", "POST").Name("Follow")
	r.HandleFunc("/api/register", handlers.API_Register).Methods("POST").Name("Follow")
	r.HandleFunc("/api/latest", handlers.API_GetLatestHandler).Methods("GET").Name("Get latest")
}

func LoadEnvVars() string {
	// Save env-vars
	dbURL := os.Getenv("DATABASE_URL")

	return dbURL
}
