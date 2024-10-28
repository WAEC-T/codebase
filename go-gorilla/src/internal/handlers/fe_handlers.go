package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-gorilla/src/internal/config"
	"go-gorilla/src/internal/db"
	"go-gorilla/src/internal/helpers"

	"html"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("SESSIONKEY"))

const PER_PAGE = 30

// Data represents the data parsed to the templates.
type Data struct {
	Messages      any
	User          any
	ProfileUser   any
	Req           string
	Followed      any
	FlashMessages []any // Changed to a slice to match the getFlash return type
}

// GetUser retrieves the user from the session.
func GetUser(r *http.Request) (any, string, error) {
	session, err := GetSession(r)
	if err != nil {
		return nil, "", err
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		return nil, "", fmt.Errorf("no user in the session")
	}

	// Perform type assertion for userID
	userIDStr := strconv.Itoa(userID.(int))
	if !ok {
		return nil, "", fmt.Errorf("user_id is not of type string")
	}

	// Query the user from the database
	user, err := db.GetUserNameByUserID(userIDStr) // Assuming queryUserByID is defined
	if err != nil {
		return nil, "", err
	}

	return user, userIDStr, nil
}

// getSession retrieves the session for the user.
func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "user-session")
}

// """Return the gravatar image for the given email address."""
func Gravatar_url(email string, size int) string {
	hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hex.EncodeToString(hash[:]), size)
}

// getFlash retrieves flash messages from the session.
func GetFlash(w http.ResponseWriter, r *http.Request) []any {
	session, err := GetSession(r)
	if err != nil {
		return nil
	}

	flashes := session.Flashes()
	session.Save(r, w)
	return flashes
}

func SetFlash(w http.ResponseWriter, r *http.Request, message string) {
	session, _ := GetSession(r)
	session.AddFlash(html.UnescapeString(message))
	session.Save(r, w)
}

func Reload(w http.ResponseWriter, r *http.Request, message string, template string) {
	d := Data{}
	if message != "" {
		SetFlash(w, r, message)
	}
	d.FlashMessages = GetFlash(w, r)
	config.Tpl.ExecuteTemplate(w, template, d)
}

// publicTimeline displays the latest messages of all users.
func Public_timeline(w http.ResponseWriter, r *http.Request) {
	user, userID, err := GetUser(r)
	if err != nil {
		// Log the error and handle the user not being logged in
		fmt.Println("public timeline: error retrieving user:", userID, err)
	}

	// Fetch public messages
	messages, err := db.GetPublicMessages(PER_PAGE)
	if err != nil {
		fmt.Println("Error fetching public messages:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare data for rendering
	flash := GetFlash(w, r)
	data := Data{
		Messages:      messages,
		User:          user,
		Req:           r.RequestURI,
		FlashMessages: flash,
		Followed:      nil,
	}

	// Render the template
	if err := config.Tpl.ExecuteTemplate(w, "timeline.html", data); err != nil {
		fmt.Println("Error trying to execute template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// """Registers the user."""
func Register(w http.ResponseWriter, r *http.Request) {
	user, _, err := GetUser(r)
	if err == nil && !(helpers.IsNil(user)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == "GET" {
		config.Tpl.ExecuteTemplate(w, "register.html", nil)

	} else if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")

		userID, err := db.GetUserIDByUsername(username)
		if err != nil {
			fmt.Errorf("failed to connect to database: %w", err)
		}

		if username == "" {
			Reload(w, r, "You have to enter a username", "register.html")
			return

		} else if !strings.Contains(email, "@") {
			Reload(w, r, "You have to enter a valid email address", "register.html")
			return

		} else if password == "" {
			Reload(w, r, "You have to enter a password", "register.html")
			return

		} else if password != password2 {
			Reload(w, r, "The two passwords do not match", "register.html")
			return

		} else if fmt.Sprint(userID) != "-1" || err != nil {
			fmt.Println("error: ", err)
			Reload(w, r, "The username is already taken", "register.html")
			return

		} else {
			hash := md5.Sum([]byte(password))
			if err != nil {
				fmt.Println("Error hashing the password")
				return
			}
			err := db.RegisterUser(username, email, hash)
			if err != nil {
				fmt.Println("error: ", err)
			}
			SetFlash(w, r, "You were successfully registered and can login now")
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

// """Logs the user in."""
func Login(w http.ResponseWriter, r *http.Request) {
	user, _, err := GetUser(r)
	if err == nil && !(helpers.IsNil(user)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == "GET" {
		Reload(w, r, "", "login.html")

	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := db.GetUserByUsername(username)
		if err != nil || helpers.IsNil(user) {
			Reload(w, r, "Invalid username", "login.html")
			return
		}

		pwHash := user.PwHash
		if !helpers.CheckPasswordHash(password, pwHash) {
			Reload(w, r, "Invalid password", "login.html")
			return
		}
		session, _ := GetSession(r)
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600, // 1 hour in seconds
			HttpOnly: true, // Recommended for security
		}
		user_id, err := db.GetUserIDByUsername(username)
		if err != nil {
			panic("This is not allowed happen!")
		}
		session.Values["user_id"] = user_id
		session.Save(r, w)
		SetFlash(w, r, "You were logged in")
		http.Redirect(w, r, "/public", http.StatusSeeOther)
		return
	}
}

// """Logs the user out"""
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		fmt.Println("Error getting session data")
	} else {
		SetFlash(w, r, "You were logged out")
		delete(session.Values, "user_id")
		err = session.Save(r, w)
		if err != nil {
			fmt.Println("Error in saving the session data")
		}
		http.Redirect(w, r, "/public", http.StatusSeeOther)
	}
}

func Timeline(w http.ResponseWriter, r *http.Request) {
	net.SplitHostPort(r.RemoteAddr)
	user, user_id, err := GetUser(r)
	if err != nil || helpers.IsNil(user) {
		http.Redirect(w, r, "/public", http.StatusFound)
	} else {

		messages, err := db.GetMyMessages(user_id)
		if err != nil {
			fmt.Println("Timeline: Error when trying to query the database", err)
			return
		}
		flash := GetFlash(w, r)
		profile_user := user

		following, err := db.GetFollowing(user_id, 30) //TODO: LIMIT OF FOLLOWERS WE QUERY?

		d := Data{
			User:          user,
			ProfileUser:   profile_user,
			Messages:      messages,
			FlashMessages: flash,
			Followed:      following,
		}

		err = config.Tpl.ExecuteTemplate(w, "timeline.html", d)
		if err != nil {
			fmt.Println("Error when trying to execute the template: ", err)
			return
		}
	}
}

// """Registers a new message for the user."""
func Add_message(w http.ResponseWriter, r *http.Request) {
	user, user_id, err := GetUser(r)
	userIDString, errStr := strconv.Atoi(user_id)
	if err != nil || errStr != nil || helpers.IsNil(user) {
		http.Error(w, "You need to login before you can post a message", http.StatusUnauthorized)
		return
	}
	text := r.FormValue("text")
	if text != "" {
		// Correct SQL query with pub_date and flagged as integer (0 for unflagged)
		err := db.AddMessage(text, userIDString)

		if err != nil {
			http.Error(w, "Unable to add message", http.StatusInternalServerError)
			return
		}

		SetFlash(w, r, "Your message was recorded")
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// """Adds the current user as follower of the given user."""
func Follow_user(w http.ResponseWriter, r *http.Request) {
	user, user_id, err := GetUser(r)
	if err != nil || helpers.IsNil(user) {
		http.Error(w, "You need to login before you can follow the user", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	println("Now following " + username)

	profileUser, err := db.GetUserByUsername(username)
	profileUserID := fmt.Sprintf("%v", profileUser.UserID)
	if err != nil {
		http.Error(w, "Followuser: Error when trying to find the user in the database in follow", http.StatusNotFound)
		return
	}
	err = db.FollowUser(user_id, profileUserID)
	if err != nil {
		fmt.Println("Error when trying to insert data into the database")
		return
	}
	message := fmt.Sprintf("You are now following &#34;%s&#34;", username)
	SetFlash(w, r, message)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// """Removes the current user as follower of the given user."""
func Unfollow_user(w http.ResponseWriter, r *http.Request) {
	user, user_id, err := GetUser(r)
	if err != nil || helpers.IsNil(user) {
		http.Error(w, "You need to login before you can follow the user", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	println("displaying username for " + username)

	profileUser, err := db.GetUserByUsername(username)
	profileUserID := fmt.Sprintf("%v", profileUser.UserID)
	if err != nil {
		http.Error(w, "Error when trying to find the user in the database in unfollow", http.StatusNotFound)
		return
	}
	err = db.UnfollowUser(user_id, profileUserID)
	if err != nil {
		fmt.Println("Error when trying to delete data from database")
		return
	}
	message := fmt.Sprintf("You are no longer following &#34;%s&#34;", username)
	SetFlash(w, r, message)
	http.Redirect(w, r, "/", http.StatusFound)
}

// """Display's a users tweets."""
func User_timeline(w http.ResponseWriter, r *http.Request) {
	user, user_id, err := GetUser(r)
	if err != nil || helpers.IsNil(user) {
		SetFlash(w, r, "You need to login before you can see the user's timeline")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]

	following, err := db.GetFollowing(user_id, 30) //TODO: LIMIT OF FOLLOWERS WE QUERY?
	profile_user, err := db.GetUserByUsername(username)
	if err != nil || helpers.IsNil(profile_user) {
		SetFlash(w, r, "The user does not exist")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	profile_user_id := profile_user.UserID

	messages, err := db.GetUserMessages(profile_user_id, 30)
	if err != nil {
		fmt.Println("User Timeline: Error when trying to query the database", err)
		return
	}
	flash := GetFlash(w, r)

	d := Data{
		Messages:      messages,
		User:          user,
		ProfileUser:   profile_user.Username,
		FlashMessages: flash,
		Followed:      following,
	}
	err = config.Tpl.ExecuteTemplate(w, "timeline.html", d)
	if err != nil {
		fmt.Println("Error when trying to execute the template: ", err)
		return
	}
}
