package model

import "database/sql"

type Club struct {
	ID           int
	Name         string
	Description  string
	DateCreated  string
	OwnerId        sql.NullInt64
	Integrations sql.NullString
	Members []User
}