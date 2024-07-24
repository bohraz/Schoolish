package model

import "database/sql"

type Post struct {
	Id       int
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	UserId   int         `json:"userId"`
	AnswerId sql.NullInt64 `json:"answerId"`
}

type Comment struct {
	Success bool 	`json:"success"`
	Id      int 	`json:"id"`
	PostId 	int 	`json:"postId"`
	User 	User 	`json:"user"`
	Content string 	`json:"content"`
}