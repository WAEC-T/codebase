package models

import "time"

type Messages struct {
	MessageID int       `gorm:"column:message_id;primaryKey"`
	AuthorID  int       `gorm:"column:author_id;not null"`
	Content   string    `gorm:"column:text;not null"`
	PubDate   time.Time `gorm:"column:pub_date;type:timestamp"`
	Flagged   int       `gorm:"column:flagged;"`
}
