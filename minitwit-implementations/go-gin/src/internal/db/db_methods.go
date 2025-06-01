package db

import (
	"fmt"
	"go-gin/src/internal/config"
	"go-gin/src/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func CheckValueInMap(maps []map[interface{}]interface{}, value interface{}) bool {
	for _, m := range maps {
		for _, v := range m {
			if v == value {
				return true
			}
		}
	}
	return false
}

func ConnectDB(uri string) (*gorm.DB, error) {
	// Create a new GORM database connection with silent logging
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Create a logger but silence it
			logger.Config{
				LogLevel: logger.Silent, // Set log level to Silent
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Fetches a username by their ID
func GetUserNameByUserID(userID int) (string, error) {
	var user models.Users
	result := config.DB.First(&user, userID) // Use the passed db instance

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return "", result.Error
	}

	return user.Username, nil
}

// fetches a user by their ID
func GetUserIDByUsername(userName string) (int, error) {
	var user models.Users
	config.DB.Where("username = ?", userName).First(&user)

	if user.UserID == 0 {
		return -1, nil
	} else {

		return user.UserID, nil
	}
}

func GetUserByUsername(userName string) (models.Users, error) {
	var user models.Users
	config.DB.Where("username = ?", userName).First(&user)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return user, config.DB.Error
	}

	return user, nil
}

func GetPublicMessages(numMsgs int) ([]models.MessageUser, error) {

	var messages []models.MessageUser
	// Ensure only the required fields are selected
	result := config.DB.Table("messages").
		Select("messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged, users.user_id, users.username, users.email").
		Joins("JOIN users ON messages.author_id = users.user_id").
		Where("messages.flagged = ?", 0).
		Order("messages.pub_date DESC").
		Limit(numMsgs).
		Find(&messages)

	if result.Error != nil {
		fmt.Println("getPublicMessages error:", config.DB.Error.Error())
		return nil, config.DB.Error
	}
	return messages, nil
}

// registers a new user
func RegisterUser(userName string, email string, password string) error {

	newUser := models.Users{
		Username: userName,
		Email:    email,
		Pwd:      password,
	}

	config.DB.Create(&newUser)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return config.DB.Error
	}

	return nil
}

// fetches all messages for the current logged in user for 'My Timeline'
func GetMyMessages(userID int) ([]models.MessageUser, []int, error) {
	var messages []models.MessageUser

	subQuery := config.DB.Table("followers").
		Select("whom_id").
		Where("who_id = ?", userID)

	var followerIDs []int

	// Find the IDs from the subquery
	if err := subQuery.Find(&followerIDs).Error; err != nil {
		fmt.Println(err.Error())
		return nil, nil, err
	}

	// Use the retrieved followerIDs in the main query
	config.DB.Table("messages").
		Select("messages.*, users.*").
		Joins("JOIN users ON messages.author_id = users.user_id").
		Where("messages.flagged = ? AND (users.user_id = ? OR users.user_id IN (?))", 0, userID, followerIDs).
		Order("messages.pub_date desc").
		Find(&messages)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return nil, nil, config.DB.Error
	}
	return messages, followerIDs, nil
}

// getFollowing fetches up to `limit` users that the user identified by userID is following
func GetFollowing(userID int, limit int) ([]models.Users, error) {
	var users []models.Users
	config.DB.
		Select("users.*").
		Joins("INNER JOIN followers ON users.user_id = followers.whom_id").
		Where("followers.who_id = ?", userID).
		Limit(limit).
		Find(&users)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return nil, config.DB.Error
	}

	return users, nil
}

func IsFollowing(userID int, visitorID int) (bool, error) {
	var follower models.Followers
	result := config.DB.
		Where("who_id = ? AND whom_id = ?", userID, visitorID).
		First(&follower)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

// adds a new message to the database
func AddMessage(text string, author_id int) error {
	currentTime, err := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error formatting time: %v", err)
	}

	newMessage := models.Messages{
		AuthorID: author_id,
		Content:  text,
		PubDate:  currentTime,
		Flagged:  0,
	}

	config.DB.Create(&newMessage)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return config.DB.Error
	}

	return nil
}

// followUser adds a new follower to the database
func FollowUser(userID int, profileUserID int) error {
	var count int64
	config.DB.Model(&models.Followers{}).Where("who_id = ? AND whom_id = ?", userID, profileUserID).Count(&count)
	if count > 0 {
		return nil
	}

	newFollower := models.Followers{
		WhoID:  userID,
		WhomID: profileUserID,
	}

	config.DB.Create(&newFollower)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return config.DB.Error
	}

	return nil
}

// unfollowUser removes a follower from the database
func UnfollowUser(userID int, profileUserID int) error {
	config.DB.Where("who_id = ? AND whom_id = ?", userID, profileUserID).Delete(&models.Followers{})

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return config.DB.Error
	}

	return nil
}

// fetches all messages from picked user
func GetUserMessages(pUserId int, numMsgs int) ([]models.MessageUser, error) {
	var messages []models.MessageUser
	config.DB.Table("messages").
		Select("messages.*, users.*").
		Joins("JOIN users ON users.user_id = messages.author_id").
		Where("users.user_id = ?", pUserId).
		Order("messages.pub_date asc").
		Limit(numMsgs).
		Find(&messages)

	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return nil, config.DB.Error
	}

	return messages, nil
}

func GetLatest() (int, error) {
	var latest models.Latest
	config.DB.Where("id = 1").First(&latest)
	return latest.Value, nil
}

func UpdateLatest(commandID int) error {
	config.DB.Save(&models.Latest{ID: 1, Value: commandID})
	if config.DB.Error != nil {
		fmt.Println(config.DB.Error.Error())
		return config.DB.Error
	}
	return nil
}
