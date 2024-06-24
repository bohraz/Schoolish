package database

import (
	"errors"
	"log"
	"root/internal/model"
	"time"
)

func GetClub(clubId int) model.Club {
	clubInfoQuery := `
		SELECT 
			clubId, 
			name, 
			description, 
			dateCreated, 
			userId, 
			integrations
		FROM app.clubs
		WHERE clubId = ?
		LIMIT 1
	`
	clubMemberListQuery := `
		SELECT username, firstName, lastName
		FROM app.users
		INNER JOIN app.users_clubs ON users.userId = users_clubs.userId
		WHERE users_clubs.clubId = ?
	`

	var club model.Club

	err := DB.QueryRow(clubInfoQuery, clubId).Scan(&club.ID, &club.Name, &club.Description, &club.DateCreated, &club.OwnerId, &club.Integrations)
	if err != nil {
		log.Println("There was an error when querying for club list: ", err)
	}

	rows, err := DB.Query(clubMemberListQuery, clubId)
	if err != nil {
		log.Println("There was an error when querying for club members: ", err)
	}

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.Handle, &user.FirstName, &user.LastName)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		club.Members = append(club.Members, user)
	}

	return club
}

func GetUserClubList(userId int) []string {
	query := `
		SELECT clubs.name
		FROM app.clubs
		INNER JOIN app.users_clubs ON clubs.clubId = users_clubs.clubId
		WHERE users_clubs.userId = ?
	`

	result, err := DB.Query(query, userId)
	if err != nil {
		log.Println("There was an error when querying for club list: ", err)
	}

	names := make([]string, 0)

	for result.Next() {
		var name string
		err = result.Scan(&name)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		names = append(names, name)
	}

	return names
}

func JoinClub(clubId, userId int) error {
	query := `
		INSERT INTO app.users_clubs (userId, clubId)
		VALUES (?, ?)
	`

	_, err := DB.Exec(query, userId, clubId)
	if err != nil {
		log.Println("There was an error joining the club:", err)
		return err
	}

	return nil
}

func LeaveClub(clubId, userId int) error {
	deletionQuery := `DELETE FROM app.users_clubs WHERE userId = ? AND clubId = ?`
	isOwnerQuery := `SELECT EXISTS(SELECT 1 FROM app.clubs WHERE userId = ? AND clubId = ?)`

	tx, err := DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else if err = tx.Commit(); err != nil {
			log.Println("There was an error committing the transaction:", err)
		}
	}()
	if err != nil {
		log.Println("There was an error starting the transaction:", err)
		return err
	}

	_, err = tx.Exec(deletionQuery, userId, clubId)
	if err != nil {
		log.Println("There was an error leaving the club:", err)
		return err
	}

	var isOwner bool
	err = tx.QueryRow(isOwnerQuery, userId, clubId).Scan(&isOwner)
	if err != nil {
		log.Println("There was an error checking if the user is the owner of the club:", err)
	} else if isOwner {
		updateQuery := `UPDATE app.clubs SET userId = NULL WHERE clubId = ?`

		_, err = DB.Exec(updateQuery, clubId)
		if err != nil {
			log.Println("There was an error updating the owner of the club", err)
			return err
		}
	}

	return nil
}

func CreateClub(name, description string, userId int) (int, error) {
	checkExistsQuery := `SELECT EXISTS(SELECT 1 FROM app.clubs WHERE name = ?)`
	clubsQuery := `INSERT INTO app.clubs (name, description, userId, dateCreated)VALUES (?, ?, ?, ?)`
	memberQuery := `INSERT INTO app.users_clubs (userId, clubId) VALUES (?, ?)`

	tx, err := DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else if err = tx.Commit(); err != nil {
			log.Println("There was an error committing the transaction:", err)
		}
	}()
	if err != nil {
		log.Println("There was an error starting the transaction:", err)
		return 0, err
	}

	var exists bool
	err = tx.QueryRow(checkExistsQuery, name).Scan(&exists)
	if err != nil {
		log.Println("There was an error checking if the club exists:", err)
		return 0, err
	}

	if exists {
		msg := "The club already exists: " + name
		log.Println(msg)
		return 0, errors.New(msg)
	}

	dateCreated := time.Now().Format("2006-01-02")
	result, err := tx.Exec(clubsQuery, name, description, userId, dateCreated)
	if err != nil {
		log.Println("There was an error creating the club:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("There was an error getting the club ID:", err)
		return 0, err
	}

	_, err = tx.Exec(memberQuery, userId, id)
	if err != nil {
		log.Println("There was an error adding the user to the club:", err)
		return 0, err
	}
	
	return int(id), nil
}