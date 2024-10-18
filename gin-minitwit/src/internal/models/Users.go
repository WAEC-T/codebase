package models

type Users struct {
	UserID   int `gorm:"primaryKey"`
	Username string
	Email    string
	PwHash   string
}
