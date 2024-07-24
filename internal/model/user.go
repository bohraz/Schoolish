package model

type User struct {
	Id             int
	Handle         string `json:"username"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
	HashedPassword string
}