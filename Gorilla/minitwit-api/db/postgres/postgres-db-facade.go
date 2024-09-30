package postgres

import (
	"errors"
	"fmt"
	"log"
	"minitwit-api/model"
	"net/url"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gorm_logger "gorm.io/gorm/logger"
)

type PostgresDbImplementation struct {
	// Implement the methods defined in the Idb interface here
	db *gorm.DB
}

func (pgImpl *PostgresDbImplementation) Connect_db() {

	user := os.Getenv("POSTGRES_USER")
	pw := os.Getenv("POSTGRES_PW")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB_NAME")

	dsn := url.URL{
		User:   url.UserPassword(user, pw),
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbname,
	}

	newLogger := gorm_logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gorm_logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)
	var err error
	pgImpl.db, err = gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("Error connecting to the database ", err)
		return
	}

	fmt.Println("Successfully connected to the database")
}

func (pgImpl *PostgresDbImplementation) QueryUserCount() float64 { // To be called each time the counters are reset (when building the image)

	var count int64
	pgImpl.db.Model(&model.Users{}).Count(&count)
	return float64(count)
}
func (pgImpl *PostgresDbImplementation) QueryMessageCount() float64 { // To be called each time the counters are reset (when building the image)

	var count int64
	pgImpl.db.Model(&model.Messages{}).Count(&count)
	return float64(count)
}
func (pgImpl *PostgresDbImplementation) QueryFollowerCount() float64 { // To be called each time the counters are reset (when building the image)

	var count int64
	pgImpl.db.Model(&model.Followers{}).Count(&count)
	return float64(count)
}

func (pgImpl *PostgresDbImplementation) QueryRegister(args []string) {
	user := &model.Users{
		Username: args[0],
		Email:    args[1],
		PwHash:   args[2],
	}
	res := pgImpl.db.Create(user)
	if res.Error != nil {
		fmt.Println("Error registering user: ", res.Error)
		return
	}
	fmt.Println("User registered successfully: ", user.Username)

}

func (pgImpl *PostgresDbImplementation) QueryMessage(message *model.Messages) {
	res := pgImpl.db.Create(message)
	if res.Error != nil {
		fmt.Println("Error creating message: ", res.Error)
		return
	}
	fmt.Println("Message created successfully: ", message.Text)
}

func (pgImpl *PostgresDbImplementation) QueryFollow(args []int) {
	followers := &model.Followers{
		WhoID:  args[0],
		WhomID: args[1],
	}
	res := pgImpl.db.Create(followers)
	if res.Error != nil {
		fmt.Println("Error creating followers: ", res.Error)
		return
	}
	fmt.Println("Follower created successfully: ", followers.WhoID, " -> ", followers.WhomID)
}

func (pgImpl *PostgresDbImplementation) QueryUnfollow(args []int) {
	res := pgImpl.db.Where("who_id = ? AND whom_id = ?", args[0], args[1]).Delete(&model.Followers{})
	if res.Error != nil {
		fmt.Println("Error unfollowing user: ", res.Error)
		return
	}
	fmt.Println("User unfollowed successfully: ", args[0], " -> ", args[1])
}

func (pgImpl *PostgresDbImplementation) QueryDelete(args []int) {
	res := pgImpl.db.Delete(&model.Users{}, args[0])
	if res.Error != nil {
		fmt.Println("Error deleting user: ", res.Error)
		return
	}
	fmt.Println("User deleted successfully: ", args[0])
}

func (pgImpl *PostgresDbImplementation) GetMessages(args []int) []map[string]any {
	var messages []model.Messages
	res := pgImpl.db.Where("flagged = 0").Order("pub_date DESC").Limit(args[0]).Find(&messages)
	if res.Error != nil {
		fmt.Println("Error getting messages: ", res.Error)
		return []map[string]any{}
	}

	var Messages []map[string]any
	for _, msg := range messages {
		var user model.Users
		pgImpl.db.First(&user, msg.AuthorID)

		message := make(map[string]any)
		message["content"] = msg.Text
		message["pub_date"] = msg.PubDate
		message["user"] = user.Username

		Messages = append(Messages, message)
	}
	fmt.Println("Messages retrieved successfully")
	return Messages
}

func (pgImpl *PostgresDbImplementation) GetMessagesForUser(args []int) []map[string]any {
	var messages []model.Messages
	res := pgImpl.db.Where("flagged = 0 AND author_id = ?", args[0]).Order("pub_date DESC").Limit(args[1]).Find(&messages)
	if res.Error != nil {
		fmt.Println("Error getting messages for user: ", res.Error)
	}

	var Messages []map[string]any

	for _, msg := range messages {
		var user model.Users
		pgImpl.db.First(&user, msg.AuthorID)

		message := make(map[string]any)
		message["content"] = msg.Text
		message["pub_date"] = msg.PubDate
		message["user"] = user.Username

		Messages = append(Messages, message)
	}
	fmt.Println("Messages for user retrieved successfully")
	return Messages
}

func (pgImpl *PostgresDbImplementation) GetFollowees(args []int) []string {
	var followees []string
	res := pgImpl.db.Model(model.Users{}).
		Select("username").
		Joins("inner join followers ON followers.whom_id = user_id").
		Where("followers.who_id = ?", args[0]).
		Limit(args[1]).
		Scan(&followees)

	if res.Error != nil {
		fmt.Println("Error getting followees: ", res.Error)
	}
	fmt.Println("Followees retrieved successfully")
	return followees
}

func (pgImpl *PostgresDbImplementation) Get_user_id(username string) (int, error) {
	var user model.Users
	res := pgImpl.db.Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("user with username '%s' not found", username)

		}
		fmt.Println("Error querying database: ", res.Error)
		return 0, fmt.Errorf("error querying database: %v", res.Error)
	}
	fmt.Println("User found: ", username)
	return user.UserID, nil
}

func (pgImpl *PostgresDbImplementation) GetAllUsers() []model.Users {
	var users []model.Users
	pgImpl.db.Find(users)
	return users
}

func (pgImpl *PostgresDbImplementation) CreateUsers(users *[]model.Users) error {
	res := pgImpl.db.CreateInBatches(&users, 100)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (pgImpl *PostgresDbImplementation) GetAllMessages() []model.Messages {
	var messages []model.Messages
	pgImpl.db.Find(messages)
	return messages
}

func (pgImpl *PostgresDbImplementation) CreateMessages(messages *[]model.Messages) error {
	res := pgImpl.db.CreateInBatches(&messages, 100)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (pgImpl *PostgresDbImplementation) GetAllFollowers() []model.Followers {
	var followers []model.Followers
	pgImpl.db.Find(followers)
	return followers
}

func (pgImpl *PostgresDbImplementation) CreateFollowers(followers *[]model.Followers) error {
	res := pgImpl.db.CreateInBatches(&followers, 100)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (pgImpl *PostgresDbImplementation) GetCount(key string) int {
	var sim model.Count
	pgImpl.db.Where("key = ?", key).First(&sim)

	return sim.Value
}

func (pgImpl *PostgresDbImplementation) SetCount(key string, value int) error {
	// Upsert operation
	upsert := pgImpl.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},              // Unique columns
		DoUpdates: clause.AssignmentColumns([]string{"value"}), // Columns to update
	}).Create(&model.Count{Key: key, Value: value})

	if upsert.Error != nil {
		log.Fatalf("failed to upsert record: %v", upsert.Error)
		return upsert.Error
	}
	return nil
}

func (pgImpl *PostgresDbImplementation) IsNil(i interface{}) bool {
	if i == nil || i == interface{}(nil) {
		return true
	} else {
		return false
	}
}

func (pgImpl *PostgresDbImplementation) IsZero(i int) bool {
	if i == 0 {
		return true
	} else {
		return false
	}
}