package database

import (
	"database/sql"
	"log"
)

// Checks if a user with given username exists and returns their password
func CheckLoginInfo(username string) string {
	query := `SELECT password FROM app.users WHERE username = ? LIMIT 1`

	var password string

	err := DB.QueryRow(query, username).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found!")
			return ""
		}
		panic(err)
	}

	return password
}

// Checks if a user with given username or email already exists
func UserFound(username, email string) (bool, string) {
	query := `SELECT username, email FROM app.users WHERE username = ? OR email = ? LIMIT 1`

	var foundFields [2]string
	err := DB.QueryRow(query, username, email).Scan(&foundFields[0], &foundFields[1])
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ""
		}
		log.Fatal(err)
	}

	if foundFields[0] == username {
		if foundFields[1] == email {
			return true, "username and email"
		}
		return true, "username"
	} else if foundFields[1] == email {
		return true, "email"
	}

	return false, "" 
}

// Inserts new user into database with given values
func CreateUser(username, email, password, firstName, lastName string) bool {
	query := `INSERT INTO app.users
			(username,email,password,firstName,lastName)
			VALUES (?,?,?,?,?)`
	
	_, err := DB.Exec(query, username, email, password, firstName, lastName)
	if err != nil {
		log.Println("There was an error executing the CreateUser query!", err)
		return false
	}

	return true
}