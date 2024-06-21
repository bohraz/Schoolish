package model

import "database/sql"

type Club struct {
	ID           int
	Name         string
	Description  string
	DateCreated  string
	Owner        sql.NullInt64
	Integrations sql.NullString
	Members []User
}