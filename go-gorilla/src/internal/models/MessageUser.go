package models

type MessageUser struct {
	MessageID int    `json:"message_id"`
	AuthorID  int    `json:"author_id"`
	Text      string `json:"content"`
	PubDate   string `json:"pub_date"`
	Flagged   int    `json:"flagged"`
	Username  string `json:"user"`
	Email     string `json:"email"`
}
