package database

import (
	"database/sql"
	"log"
	"root/internal/model"
)

// Checks if a user with given username exists and returns their password
func CheckLoginInfo(username string) (model.User, string) {
	query := `SELECT userId, password, firstName, lastName, email FROM app.users WHERE username = ? LIMIT 1`

	var password string
	user := model.User{
		Handle: username,
		// Email: request.FormValue("email"),
		// FirstName: request.FormValue("firstname"),
		// LastName: request.FormValue("lastname"),
	}

	err := DB.QueryRow(query, username).Scan(&user.Id, &password, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found!")
			return user, ""
		}
		panic(err)
	}

	return user, password
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
func CreateUser(username, email, password, firstName, lastName string) (uint, error) {
	query := `INSERT INTO app.users
			(username,email,password,firstName,lastName)
			VALUES (?,?,?,?,?)`
	
	result, err := DB.Exec(query, username, email, password, firstName, lastName)
	if err != nil {
		log.Println("There was an error executing the CreateUser query!", err)
		return 0, err
	}

	userId, err := result.LastInsertId()
    if err != nil {
        log.Println("There was an error retrieving the last insert ID!", err)
        return 0, err
    }

	return uint(userId), nil
}