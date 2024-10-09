package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func format_datetime(timestamp string) string {
	// Define the layout that matches your input timestamp format
	layout := time.RFC3339 // or "2006-01-02T15:04:05Z07:00" for a custom layout

	// Parse the string into a time.Time object
	parsedTime, err := time.Parse(layout, timestamp)
	if err != nil {
		// Handle the error, return the original string or an error message
		fmt.Println("Error parsing timestamp:", err)
		return timestamp
	}

	// Format the time.Time object into your desired display format
	return parsedTime.Format("2006-01-02 15:04:05") // Customize this layout as needed
}

return i == nil || i == interface{}(nil)
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

func checkPasswordHash(userEnteredPwd string, dbpwd string) bool {
	hash := md5.Sum([]byte(userEnteredPwd))
	str := hex.EncodeToString(hash[:])
	return str == dbpwd
}
