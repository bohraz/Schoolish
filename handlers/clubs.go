package handlers

import (
	"fmt"
	"log"
	"net/http"
	"root/internal/database"
	"root/internal/model"
	"strconv"
	"strings"
)

func getClubFromURL(writer http.ResponseWriter, request *http.Request, idIndex int) (model.Club) {
	clubId, err := strconv.Atoi(strings.Split(request.URL.Path, "/")[idIndex])
	if err != nil {
		http.Error(writer, "There was an error processing the URL", http.StatusBadRequest)
		return model.Club{}
	}

	club := database.GetClub(clubId)
	if club.ID == 0 {
		http.Error(writer, "Invalid club ID", http.StatusBadRequest)
		return club
	}

	return club
}

func Club(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	
	clubList := database.GetUserClubList(user.Id)
	club := getClubFromURL(writer, request, 2)
	
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
	user, _ := GetUser(request)
	club := getClubFromURL(writer, request, 3)

	err := database.JoinClub(club.ID, user.Id)
	if err != nil {
		http.Error(writer, "Error joining club", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf(`<div id="message">You have joined %s!</div>`, club.Name)
	script := fmt.Sprintf(`<script>setTimeout(function() { window.location.href = "/clubs/%d/"; }, 2000);</script>`, club.ID)

	fmt.Fprintln(writer, message)
	fmt.Fprintln(writer, script)
}

func ClubLeave(writer http.ResponseWriter, request *http.Request) {
	user, _ := GetUser(request)
	club := getClubFromURL(writer, request, 3)

	err := database.LeaveClub(club.ID, user.Id)
	if err != nil {
		http.Error(writer, "Error leaving club", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf(`<div id="message">You have left %s!</div>`, club.Name)
	script := fmt.Sprintf(`<script>setTimeout(function() { window.location.href = "/clubs/%d/"; }, 2000);</script>`, club.ID)

	fmt.Fprintln(writer, message)
	fmt.Fprintln(writer, script)
}

func ClubSearch(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Club search!")
}

func ClubCreateSubmit(writer http.ResponseWriter, request *http.Request) {
	var (
		name = request.FormValue("clubName")
		desc = request.FormValue("clubDescription")
	)

	user, err := GetUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	clubId, err := database.CreateClub(name, desc, user.Id)
	if err != nil {
		http.Error(writer, "Error creating club", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(writer, `<div id="message">Club %s created!</div>`, name)
	fmt.Fprintf(writer, `<script>setTimeout(function() { window.location.href = "/clubs/%d/"; }, 2000);</script>`, clubId)
}