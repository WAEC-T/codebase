package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gin-minitwit/src/internal/models"
	"os"
	"reflect"
	"strings"
)

// Helper functions

func CheckPasswordHash(userEnteredPwd string, dbpwd string) bool {
	hash := md5.Sum([]byte(userEnteredPwd))
	str := hex.EncodeToString(hash[:])
	fmt.Println("test")
	return str == dbpwd
}

func GravatarURL(email string, size int) string {
	if size <= 0 {
		size = 80 // Default size
	}

	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%x?d=identicon&s=%d", hash, size)
}

func FormatMessages(messages []models.MessageUser) []models.MessageUI {
	var formattedMessages []models.MessageUI

	for _, m := range messages {
		var msg models.MessageUI
		// Use type assertion for int64, then convert to int
		if reflect.TypeOf(m.MessageID).Kind() == reflect.Int {
			msg.MessageID = m.MessageID
		}
		if reflect.TypeOf(m.AuthorID).Kind() == reflect.Int {
			msg.AuthorID = m.AuthorID
		}
		if reflect.TypeOf(m.UserID).Kind() == reflect.Int {
			msg.User.UserID = m.UserID
		}
		if reflect.TypeOf(m.Text).Kind() == reflect.String {
			msg.Text = m.Text
		}
		if reflect.TypeOf(m.Username).Kind() == reflect.String {
			msg.Username = m.Username
		}
		if reflect.TypeOf(m.Email).Kind() == reflect.String {
			msg.Email = m.Email
		}
		if reflect.TypeOf(m.PubDate).Kind() == reflect.Int {
			msg.PubDate = m.PubDate
		}
		link := "/" + msg.Username
		msg.Profile_link = strings.ReplaceAll(link, " ", "%20")

		gravatarURL := GravatarURL(msg.Email, 48)
		msg.Gravatar = gravatarURL

		formattedMessages = append(formattedMessages, msg)
	}

	return formattedMessages
}

func FilterMessages(messages []models.MessageUser) []models.FilteredMsg {
	var filteredMessages []models.FilteredMsg
	for _, m := range messages {
		var filteredMsg models.FilteredMsg
		// content
		if reflect.TypeOf(m.Text).Kind() == reflect.String {
			filteredMsg.Content = m.Text
		}

		// publication date
		filteredMsg.PubDate = int64(m.PubDate)

		// user
		if reflect.TypeOf(m.Username).Kind() == reflect.String {
			filteredMsg.User = m.Username
		}

		filteredMessages = append(filteredMessages, filteredMsg)
	}
	return filteredMessages
}

func LogMessage(message string) {
	// Specify the file path
	filePath := "./tmp/logging/logger.txt"

	// Open or create the file for writing
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	data := []byte(message + "\n")

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
