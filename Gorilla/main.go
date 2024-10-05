package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	postgresDB *gorm.DB
	tpl        *template.Template // Declare tpl here for global access
)

func main() {
	var err error

	/*---------------------
	* Load env vars
	*----------------------*/
	dsn := "user=user password=password dbname=waect host=localhost port=5432 sslmode=disable" // Modify as needed

	/*---------------------
	 * Connect to DB
	 *----------------------*/
	postgresDB, err = connectDB(dsn) // Ensure this matches the function name in the db package
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	/*---------------------
	* Setup routing
	*----------------------*/
	funcMap := template.FuncMap{
		"getavatar": func(url string, size int) string {
			return gravatar_url(url, size) // Assuming this function is defined somewhere
		},
		"gettimestamp": func(time string) string {
			return format_datetime(time) // Assuming this function is defined somewhere
		},
		"url_for": func(routename string, username string) string {
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
			return checkValueInMap(following, messageAuthorId)
		},
	}

	tpl, err = template.New("timeline.html").Funcs(funcMap).ParseGlob("templates/*.html") // We need to add the funcs that we want to use before parsing
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/public", public_timeline)
	r.HandleFunc("/register", register)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/", timeline)
	r.HandleFunc("/add_message", add_message).Methods("POST")
	r.HandleFunc("/{username}/follow", follow_user)
	r.HandleFunc("/{username}", user_timeline)
	r.HandleFunc("/{username}/unfollow", unfollow_user)

	fmt.Println("Listening on port 5000...")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
