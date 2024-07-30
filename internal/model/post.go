package model

import "database/sql"

var PostsBroadcast = make(chan Post)

type Post struct {
	Id       int
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	UserId   int         `json:"userId"`
	AnswerId sql.NullInt64 `json:"answerId"`
	Comments int 	   `json:"comments"`
}

type Comment struct {
	Success bool 	`json:"success"`
	Id      int 	`json:"id"`
	PostId 	int 	`json:"postId"`
	User 	User 	`json:"user"`
	Content string 	`json:"content"`
	ReplyId sql.NullInt64 `json:"replyId"`
}