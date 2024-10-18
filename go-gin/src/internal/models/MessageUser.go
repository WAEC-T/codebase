package models

import "time"

type MessageUser struct {
	MessageID int `gorm:"primaryKey"`
	AuthorID  int
	Text      string
	PubDate   time.Time
	Flagged   int
	UserID    int `gorm:"primaryKey"`
	Username  string
	Email     string
	PwHash    string
}
