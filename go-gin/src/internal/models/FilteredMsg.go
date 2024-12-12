package models

type FilteredMsg struct {
	Content string `json:"content"`
	PubDate int64  `json:"pub_date"`
	User    string `json:"user"`
}
