package model

import "database/sql"

type Club struct {
	ID           int
	Name         string
	Description  string
	DateCreated  string
	Owner        int
	Integrations sql.NullString
	Members []User
}