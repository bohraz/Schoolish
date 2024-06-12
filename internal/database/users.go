package database

import (
	"database/sql"
)

func CheckLoginInfo(username, password string) bool {
	query := `SELECT 1 FROM users WHERE username = ? AND password = ? LIMIT 1`

	var exists bool
	// Checks if a user with given username and password exists
	err := DB.QueryRow(query, username, password).Scan(&exists)
	if err != nil {
		// Return false if no rows were found
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	}

	return true
}