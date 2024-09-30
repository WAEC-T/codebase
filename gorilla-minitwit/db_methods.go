package main

import (
	"encoding/hex"
	"fmt"
	"gorilla-minitwit/models"

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

	fmt.Println("messages: ", messages)

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
