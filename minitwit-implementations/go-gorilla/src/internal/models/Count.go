package models

type Count struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value int    `json:"value"`
}
