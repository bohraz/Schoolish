package database

import (
	"database/sql"
	"fmt"
	"log"
	"root/internal/auth"
	"root/internal/model"
)

func QueryUser(identifier interface{}) (model.User, error) {
	query := `SELECT userId, username, password, firstName, lastName, email FROM app.users WHERE username = ? OR userId = ? LIMIT 1`

	user := model.User{}

	err := DB.QueryRow(query, identifier, identifier).Scan(&user.Id, &user.Handle, &user.HashedPassword, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		panic(err)
	}

	return user, nil
}

// Checks if a user with given username exists and returns their password
func GetLoginInfo(username string) (int, string) {
	user, err := QueryUser(username)
	if err != nil {
		log.Println("There was an error querying the user!", err)
		return 0, ""
	}

	return user.Id, user.HashedPassword
}

// Checks if a user with given username or email already exists and returns what, if anything, exists in a string
func UserFound(username, email string) (bool, string) {
	user, err := QueryUser(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ""
		}
		log.Println("There was an error querying the user!", err)
		return false, ""
	}

	if user.Handle == username {
		if user.Email == email {
			return true, "username and email"
		} else {
			return true, "username"
		}
	} else if user.Email == email {
		return true, "email"
	}

	return false, "" 
}

// Inserts new user into database with given values
func CreateUser(username, email, password, firstName, lastName string) (int, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		fmt.Println("There was an error hashing the password!", err)
	}

	query := `INSERT INTO app.users
			(username,email,password,firstName,lastName)
			VALUES (?,?,?,?,?)`
	
	result, err := DB.Exec(query, username, email, hashedPassword, firstName, lastName)
	if err != nil {
		log.Println("There was an error executing the CreateUser query!", err)
		return 0, err
	}

	userId, err := result.LastInsertId()
    if err != nil {
        log.Println("There was an error retrieving the last insert ID!", err)
        return 0, err
    }

	return int(userId), nil
}