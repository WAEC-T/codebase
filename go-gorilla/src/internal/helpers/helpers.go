package helpers

import (
	"fmt"
	"time"
)

func Format_datetime(timestamp string) string {
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
	return parsedTime.Format("2006-01-02 @ 15:04") // Customize this layout as needed
}

func IsNil(i interface{}) bool {
	return i == nil || i == interface{}(nil)
}

func CheckPassword(userEnteredPwd string, dbpwd string) bool {
	return userEnteredPwd == dbpwd
}
