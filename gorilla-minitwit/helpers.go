package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// """Format a timestamp for display."""
func format_datetime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 @ 15:04")
}

func isNil(i interface{}) bool {
	if i == nil || i == interface{}(nil) {
		return true
	} else {
		return false
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
