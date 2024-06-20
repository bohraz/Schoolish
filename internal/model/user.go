package model

type User struct {
	Id             uint
	Handle         string
	FirstName      string
	LastName       string
	HashedPassword string
	Email          string
}