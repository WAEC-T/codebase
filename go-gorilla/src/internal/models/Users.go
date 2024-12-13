package models

type Users struct {
	UserID   int    `gorm:"column:user_id;primaryKey"`
	Username string `gorm:"column:username;not null;index"`
	Email    string `gorm:"column:email;not null"`
	Pwd      string `gorm:"column:pw_hash;not null"`
}

func (Users) TableName() string {
	return "users"
}
