package models

type Followers struct {
	WhoID  int `gorm:"column:who_id;index:idx_member;primaryKey"`
	WhomID int `gorm:"column:whom_id;index:idx_member;primaryKey"`
}
