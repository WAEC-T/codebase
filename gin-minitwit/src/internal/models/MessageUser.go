package models

type MessageUser struct {
	MessageID int `gorm:"primaryKey"`
	AuthorID  int
	Text      string
	PubDate   int
	Flagged   int
	UserID    int `gorm:"primaryKey"`
	Username  string
	Email     string
	PwHash    string
}
