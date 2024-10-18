package models

import "time"

type MessageUI struct {
	MessageID    int
	AuthorID     int
	Text         string
	PubDate      time.Time
	Flagged      bool
	User         Users
	Email        string
	Username     string
	Profile_link string
	Gravatar     string
}
