package handlers

import (
	"fmt"
	"net/http"
	"root/internal/database"
	"root/internal/model"
	"strconv"
	"strings"
)

func getClubFromURL(request *http.Request) (model.Club, error) {
	clubId, err := strconv.Atoi(strings.Split(request.URL.Path, "/")[2])
	if err != nil {
		return model.Club{}, err
	}

	club := database.GetClub(uint(clubId))

	return club, nil
}


func ClubSearch(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Club search!")
}

func Club(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
	
	clubList := database.GetUserClubList(user.Id)
	club, err := getClubFromURL(request)
	if err != nil {
		http.Error(writer, "Invalid club ID", http.StatusBadRequest)
	}
	
	fmt.Fprintln(writer, "Selected Club:", club.ID, club.Name, club.Description, club.Owner, club.DateCreated)
	
	memberCount := len(club.Members)
	fmt.Fprint(writer, "Members:")
	for i, v := range club.Members {
		if memberCount == i + 1 {
			fmt.Fprintln(writer, v.FirstName, v.LastName + ".")
		} else {
			fmt.Fprint(writer, v.FirstName, v.LastName + ", ")
		}
	}

	for i, v := range clubList {
		fString := fmt.Sprintf("Club %v: %s", i + 1, v)
		fmt.Fprintln(writer, fString)
	}
}

func ClubJoin(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Club join!")
}

func ClubCreate(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Club create!")
}

func ClubCreateSubmit(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Club creation submit!")
}