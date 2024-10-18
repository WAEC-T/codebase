package models

type MessageUser struct {
	MessageID int `gorm:"primaryKey"`
	AuthorID  int
	Text      string
	Pub_date  string
	Flagged   int
	UserID    int `gorm:"primaryKey"`
	Username  string
	Email     string
}
