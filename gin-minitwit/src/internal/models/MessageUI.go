package models

type MessageUI struct {
	MessageID    int
	AuthorID     int
	Text         string
	PubDate      int
	Flagged      bool
	User         Users
	Email        string
	Username     string
	Profile_link string
	Gravatar     string
}
