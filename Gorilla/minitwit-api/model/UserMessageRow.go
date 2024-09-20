package model

type UserMessageRow struct {
	Message_id int
	Author_id  int
	Text       string
	Pub_date   int
	Flagged    int
	User_id    int
	Username   string
	Email      string
	Pw_hash    string
}
