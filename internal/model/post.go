package model

type Post struct {
	Id      int
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  int    `json:"userId"`
}