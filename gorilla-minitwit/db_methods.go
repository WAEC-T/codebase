package main

import (
	"encoding/hex"
	"fmt"
	"gorilla-minitwit/models"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func checkValueInMap(maps []map[interface{}]interface{}, value interface{}) bool {
	for _, m := range maps {
		for _, v := range m {
			if v == value {
				return true
			}
		}
	}
	return false
}

func connectDB(dsn string) (*gorm.DB, error) {
	fmt.Println("Connecting to the database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Perform Auto Migration
	if err := db.AutoMigrate(&models.Users{}, &models.Messages{}, &models.Followers{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Fetches a username by their ID
func getUserNameByUserID(userID string) (string, error) {
	var user models.Users
	result := postgresDB.First(&user, userID) // Use the passed db instance

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return "", result.Error
	}

	return user.Username, nil
}

// fetches a user by their ID
func getUserIDByUsername(userName string) (int, error) {
	var user models.Users
	postgresDB.Where("username = ?", userName).First(&user)

	if user.UserID == 0 {
		return -1, nil
	} else {

		return user.UserID, nil
	}
}

func getUserByUsername(userName string) (models.Users, error) {
	var user models.Users
	postgresDB.Where("username = ?", userName).First(&user)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return user, postgresDB.Error
	}

	return user, nil
}

func getPublicMessages(numMsgs int) ([]models.MessageUser, error) {

	var messages []models.MessageUser
	// Ensure only the required fields are selected
	result := postgresDB.Table("messages").
		Select("messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged, users.user_id, users.username, users.email").
		Joins("JOIN users ON messages.author_id = users.user_id").
		Where("messages.flagged = ?", 0).
		Order("messages.pub_date DESC").
		Limit(numMsgs).
		Find(&messages)

	if result.Error != nil {
		fmt.Println("getPublicMessages error:", postgresDB.Error.Error())
		return nil, postgresDB.Error
	}
	return messages, nil
}

// registers a new user
func registerUser(userName string, email string, password [16]byte) error {

	pwHashString := hex.EncodeToString(password[:])

	newUser := models.Users{
		Username: userName,
		Email:    email,
		PwHash:   pwHashString,
	}

	postgresDB.Create(&newUser)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return postgresDB.Error
	}

	return nil
}

// fetches all messages for the current logged in user for 'My Timeline'
func getMyMessages(userID string) ([]models.MessageUser, error) {

	var messages []models.MessageUser

	subQuery := postgresDB.Table("followers").
		Select("whom_id").
		Where("who_id = ?", userID)

	var followerIDs []int

	// Find the IDs from the subquery
	if err := subQuery.Find(&followerIDs).Error; err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Use the retrieved followerIDs in the main query
	postgresDB.Table("messages").
		Select("messages.*, users.*").
		Joins("JOIN users ON messages.author_id = users.user_id").
		Where("messages.flagged = ? AND (users.user_id = ? OR users.user_id IN (?))", 0, userID, followerIDs).
		Order("messages.pub_date desc").
		Find(&messages)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return nil, postgresDB.Error
	}
	return messages, nil
}

// getFollowing fetches up to `limit` users that the user identified by userID is following
func getFollowing(userID string, limit int) ([]map[interface{}]interface{}, error) {
	var users []models.Users
	postgresDB.
		Select("users.*").
		Joins("INNER JOIN followers ON users.user_id = followers.whom_id").
		Where("followers.who_id = ?", userID).
		Limit(limit).
		Find(&users)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return nil, postgresDB.Error
	}

	// Convert []models.Users to []map[interface{}]interface{}
	var result []map[interface{}]interface{}
	for _, user := range users {
		m := map[interface{}]interface{}{
			"UserID": user.UserID,
			"Name":   user.Username,
			"Email":  user.Email,
			// Add more fields as needed
		}
		result = append(result, m)
	}

	return result, nil
}

// adds a new message to the database
func addMessage(text string, author_id int) error {
	currentTime := time.Now().UTC()
	unixTimestamp := currentTime.Unix()

	newMessage := models.Messages{
		AuthorID: author_id,
		Text:     text,
		PubDate:  string(unixTimestamp), //TODO: ALIGN W. LADS: IS THIS CORRECT?
		Flagged:  0,
	}

	postgresDB.Create(&newMessage)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return postgresDB.Error
	}

	return nil
}

// followUser adds a new follower to the database
func followUser(userID string, profileUserID string) error {

	userIDInt, errz := strconv.Atoi(userID)
	profileUserIDInt, errx := strconv.Atoi(profileUserID)

	if errz != nil {
		fmt.Println(errz.Error())
		return errz
	} else if errx != nil {
		fmt.Println(errx.Error())
		return errx
	}

	// following relationship already exists
	var count int64
	postgresDB.Model(&models.Followers{}).Where("who_id = ? AND whom_id = ?", userIDInt, profileUserIDInt).Count(&count)
	if count > 0 {
		return nil
	}

	newFollower := models.Followers{
		WhoID:  userIDInt,
		WhomID: profileUserIDInt,
	}

	postgresDB.Create(&newFollower)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return postgresDB.Error
	}

	return nil
}

// unfollowUser removes a follower from the database
func unfollowUser(userID string, profileUserID string) error {
	userIDInt, errz := strconv.Atoi(userID)
	profileUserIDInt, errx := strconv.Atoi(profileUserID)

	if errz != nil {
		fmt.Println(errz.Error())
		return errz
	} else if errx != nil {
		fmt.Println(errx.Error())
		return errx
	}

	postgresDB.Where("who_id = ? AND whom_id = ?", userIDInt, profileUserIDInt).Delete(&models.Followers{})

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return postgresDB.Error
	}

	return nil
}

// fetches all messages from picked user
func getUserMessages(pUserId int, numMsgs int) ([]models.MessageUser, error) {
	var messages []models.MessageUser
	postgresDB.Table("messages").
		Select("messages.*, users.*").
		Joins("JOIN users ON users.user_id = messages.author_id").
		Where("users.user_id = ?", pUserId).
		Order("messages.pub_date asc").
		Limit(numMsgs).
		Find(&messages)

	if postgresDB.Error != nil {
		fmt.Println(postgresDB.Error.Error())
		return nil, postgresDB.Error
	}

	return messages, nil
}
