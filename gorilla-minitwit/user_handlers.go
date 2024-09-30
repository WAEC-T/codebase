package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

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

// getUser retrieves the user from the session.
func getUser(r *http.Request) (any, string, error) {
	session, err := getSession(r)
	if err != nil {
		return nil, "", err
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		return nil, "", fmt.Errorf("no user in the session")
	}

	fmt.Println("getUser user_id:", userID)

	// Perform type assertion for userID
	userIDStr := strconv.Itoa(userID.(int))
	if !ok {
		return nil, "", fmt.Errorf("user_id is not of type string")
	}

	// Query the user from the database
	user, err := getUserNameByUserID(userIDStr) // Assuming queryUserByID is defined
	if err != nil {
		return nil, "", err
	}

	return user, userIDStr, nil
}

// getSession retrieves the session for the user.
func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "user-session")
}

// """Return the gravatar image for the given email address."""
func gravatar_url(email string, size int) string {
	hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hex.EncodeToString(hash[:]), size)
}

// getFlash retrieves flash messages from the session.
func getFlash(w http.ResponseWriter, r *http.Request) []any {
	session, err := getSession(r)
	if err != nil {
		return nil
	}

	flashes := session.Flashes()
	session.Save(r, w)
	return flashes
}

func setFlash(w http.ResponseWriter, r *http.Request, message string) {
	session, _ := getSession(r)
	session.AddFlash(html.UnescapeString(message))
	session.Save(r, w)
}

func reload(w http.ResponseWriter, r *http.Request, message string, template string) {
	d := Data{}
	if message != "" {
		setFlash(w, r, message)
	}
	d.FlashMessages = getFlash(w, r)
	tpl.ExecuteTemplate(w, template, d)
}

// publicTimeline displays the latest messages of all users.
func public_timeline(w http.ResponseWriter, r *http.Request) {
	user, userID, err := getUser(r)
	fmt.Println("user: ", user)
	fmt.Println("user_id:", userID)
	if err != nil {
		// Log the error and handle the user not being logged in
		fmt.Println("public timeline: error retrieving user:", err)
	}

	// Fetch public messages
	messages, err := getPublicMessages(PER_PAGE)
	fmt.Println("getPublicMessages messages:", messages)
	if err != nil {
		fmt.Println("Error fetching public messages:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare data for rendering
	flash := getFlash(w, r)
	data := Data{
		Messages:      messages,
		User:          user,
		Req:           r.RequestURI,
		FlashMessages: flash,
		Followed:      nil,
	}

	// Render the template
	if err := tpl.ExecuteTemplate(w, "timeline.html", data); err != nil {
		fmt.Println("Error trying to execute template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// """Registers the user."""
func register(w http.ResponseWriter, r *http.Request) {
	user, _, err := getUser(r)
	if err == nil && !(isNil(user)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "register.html", nil)

	} else if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")

		if username == "" {
			reload(w, r, "You have to enter a username", "register.html")
			return

		} else if !strings.Contains(email, "@") {
			reload(w, r, "You have to enter a valid email address", "register.html")
			return

		} else if password == "" {
			reload(w, r, "You have to enter a password", "register.html")
			return

		} else if password != password2 {
			reload(w, r, "The two passwords do not match", "register.html")
			return

		} else if id, err := getUserIDByUsername(username); err != nil && id != 0 {
			fmt.Println("error: ", err)
			reload(w, r, "The username is already taken", "register.html")
			return

		} else {
			hash := md5.Sum([]byte(password))
			if err != nil {
				fmt.Println("Error hashing the password")
				return
			}
			err := registerUser(username, email, hash)
			if err != nil {
				fmt.Println("error: ", err)
			}
			setFlash(w, r, "You were successfully registered and can login now")
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

// """Logs the user in."""
func login(w http.ResponseWriter, r *http.Request) {
	user, _, err := getUser(r)
	if err == nil && !(isNil(user)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == "GET" {
		reload(w, r, "", "login.html")

	} else if r.Method == "POST" {
		fmt.Println("POST, render login")
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := getUserByUsername(username)
		if err != nil || isNil(user) {
			reload(w, r, "Invalid username", "login.html")
			return
		}

		pwHash := user.PwHash
		if !checkPasswordHash(password, pwHash) {
			reload(w, r, "Invalid password", "login.html")
			return
		}
		session, _ := getSession(r)
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600, // 1 hour in seconds
			HttpOnly: true, // Recommended for security
		}
		user_id, err := getUserIDByUsername(username)
		if err != nil {
			panic("This is not allowed happen!")
		}
		session.Values["user_id"] = user_id
		fmt.Println("setting user_id for session: ", session.Values["user_id"])
		session.Save(r, w)
		setFlash(w, r, "You were logged in")
		http.Redirect(w, r, "/public", http.StatusSeeOther)
		return
	}
}
