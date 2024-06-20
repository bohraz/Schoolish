package model

type User struct {
	Id             int
	Handle         string
	FirstName      string
	LastName       string
	HashedPassword string
	Email          string
}