package model

type User struct {
	Id             int
	Handle         string `json:"username"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Password       string `json:"password"`
	HashedPassword string // Thou shalt not json marshal this field for security reasons
}