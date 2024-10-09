package main

import (
	"fmt"
	"gorilla-minitwit/internal/config"
	"gorilla-minitwit/internal/db"
	"gorilla-minitwit/internal/handlers"
	"gorilla-minitwit/internal/helpers"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var err error

func main() {

	/*---------------------
	 * Load env vars
	 *----------------------*/
	dsn := LoadEnvVars()

	/*---------------------
	 * Connect to DB
	 *----------------------*/
	config.DB, err = db.ConnectDB(dsn) // Ensure this matches the function name in the db package
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	/*---------------------
	* Setup routing
	*----------------------*/
	funcMap := SetupRouting()

	config.Tpl, err = template.New("timeline.html").Funcs(funcMap).ParseGlob("../../templates/*.html") // We need to add the funcs that we want to use before parsing
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	/*---------------------
	* Setup route-handlers
	*----------------------*/
	r := mux.NewRouter()
	SetRouteHandlers(r)

	fmt.Println("Listening on port 5000...")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

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
			case "timeline":
				return "/" + username
			case "public_timeline":
				return "/public"
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
		"IsFollowing": func(following []map[interface{}]interface{}, messageAuthorId int) bool {
			return db.CheckValueInMap(following, messageAuthorId)
		},
	}
	return funcMap
}

func SetRouteHandlers(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/public", handlers.Public_timeline)
	r.HandleFunc("/register", handlers.Register)
	r.HandleFunc("/login", handlers.Login)
	r.HandleFunc("/logout", handlers.Logout)
	r.HandleFunc("/", handlers.Timeline)
	r.HandleFunc("/add_message", handlers.Add_message).Methods("POST")
	r.HandleFunc("/{username}/follow", handlers.Follow_user)
	r.HandleFunc("/{username}", handlers.User_timeline)
	r.HandleFunc("/{username}/unfollow", handlers.Unfollow_user)
}

func LoadEnvVars() string {
	// Load the .env file
	err = godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Save env-vars
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("user=%s password=%s dbname=waect host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort)
	fmt.Println("DSN:", dsn)

	return dsn
}
