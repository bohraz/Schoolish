package database

import (
	"log"
	"root/internal/model"
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

	err := DB.QueryRow(clubInfoQuery, clubId).Scan(&club.ID, &club.Name, &club.Description, &club.DateCreated, &club.Owner, &club.Integrations)
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
	query := `
		DELETE FROM app.users_clubs
		WHERE userId = ? AND clubId = ?
	`

	_, err := DB.Exec(query, userId, clubId)
	if err != nil {
		log.Println("There was an error leaving the club:", err)
		return err
	}

	return nil
}

func CreateClub(userId, clubId int, name string) {

}