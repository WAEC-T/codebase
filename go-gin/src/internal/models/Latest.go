package models

type Latest struct {
	ID    int `gorm:"primaryKey"`
	Value int
}
