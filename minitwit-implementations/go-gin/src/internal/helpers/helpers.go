package helpers

import (
	"crypto/md5"
	"fmt"
	"go-gin/src/internal/models"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Helper functions
func GravatarURL(email string, size int) string {
	if size <= 0 {
		size = 80 // Default size
	}

	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%x?d=identicon&s=%d", hash, size)
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
		msg.PubDate = Format_datetime(m.PubDate) // Assuming PubDate is already a time.Time type

		// Create the profile link by encoding the username
		msg.Profile_link = "/user/" + strings.ReplaceAll(msg.Username, " ", "%20")

		// Generate the Gravatar URL
		msg.Gravatar = GravatarURL(msg.Email, 48)

		// Append the formatted message to the result slice
		formattedMessages = append(formattedMessages, msg)
	}

	return formattedMessages
}

func Format_datetime(timestamp time.Time) string {

	// Format the time.Time object into your desired display format
	return timestamp.Format("2006-01-02 @ 15:04") // Customize this layout as needed
}

func SaveSessionOrRedirect(c *gin.Context, err error, redirectURL string) bool {
	if err != nil {
		fmt.Println("session save failed with:", err)
		c.Redirect(http.StatusFound, redirectURL)
		return false
	}
	return true
}

func CheckPassword(userEnteredPwd string, dbpwd string) bool {
	return userEnteredPwd == dbpwd
}

func IsNil(i interface{}) bool {
	return i == nil || i == interface{}(nil)
}
