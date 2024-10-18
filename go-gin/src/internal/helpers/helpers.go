package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-gin/src/internal/models"
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
		// Initialize a new MessageUI struct
		var msg models.MessageUI

		// Directly assign fields from MessageUser to MessageUI
		msg.MessageID = m.MessageID
		msg.AuthorID = m.AuthorID
		msg.User.UserID = m.UserID
		msg.Text = m.Text
		msg.Username = m.Username
		msg.Email = m.Email
		msg.PubDate = m.PubDate // Assuming PubDate is already a time.Time type

		// Create the profile link by encoding the username
		msg.Profile_link = "/" + strings.ReplaceAll(msg.Username, " ", "%20")

		// Generate the Gravatar URL
		msg.Gravatar = GravatarURL(msg.Email, 48)

		// Append the formatted message to the result slice
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
