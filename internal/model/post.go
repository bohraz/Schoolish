package model

import "database/sql"

type Post struct {
	Id       int
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	UserId   int         `json:"userId"`
	AnswerId sql.NullInt64 `json:"answerId"`
}