package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"root/internal/database"
	"root/internal/model"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var ErrNotAClubUrl = errors.New("this is not a club url or there is no clubid")

// Get the club ID from the URL and check if the user has permission to edit the club if requesting an edit page
func getClubIdFromUrl(writer http.ResponseWriter, request *http.Request) (int, error) {
	isClub := false
	wantEdit := false

	splitPath := strings.Split(request.URL.Path, "/")
	removeEmptyStringsFromPath(&splitPath)

	for i, v := range splitPath {
		if v == "clubs" || v == "club" {
			isClub = true
		} else if v == "edit" {
			wantEdit = true
		} else if isClub && (i == len(splitPath) - 1) {
			clubId, err := strconv.Atoi(v)
			if err != nil {
				http.Error(writer, "Invalid club ID", http.StatusBadRequest)
			}

			if wantEdit {
				user, err := GetLoggedInUser(request)
				if err != nil {
					http.Error(writer, "Unauthorized", http.StatusUnauthorized)
					return 0, err
				}
				role := database.GetUserClubRole(user.Id, clubId)
				if role < 2 {
					http.Error(writer, "Unauthorized", http.StatusUnauthorized)
					return 0, err
				}
			}

			return clubId, nil
		}
	}

	return 0, ErrNotAClubUrl
}

func getClubFromUrl(writer http.ResponseWriter, request *http.Request) (model.Club) {
	clubIdStr := mux.Vars(request)["id"]
	clubId, err := strconv.Atoi(clubIdStr)
	if err != nil {
		http.Error(writer, "Invalid club ID", http.StatusBadRequest)
		return model.Club{}
	}

	club := database.GetClub(clubId)
	if club.ID == 0 {
		http.Error(writer, "Invalid club ID", http.StatusBadRequest)
		return club
	}

	return club
}

func CreateClubApi(writer http.ResponseWriter, request *http.Request) {
	user, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	var createRequest model.Club
	err = json.NewDecoder(request.Body).Decode(&createRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}


	clubId, err := database.CreateClub(createRequest.Name, createRequest.Description, user.Id)
	if err != nil {
		http.Error(writer, "Error creating club", http.StatusInternalServerError)
		return
	}

	response := model.Club{
		ID: clubId,
		Name: createRequest.Name,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJson)
}

func EditClubApi(writer http.ResponseWriter, request *http.Request) {
	clubId, err := getClubIdFromUrl(writer, request)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedClub := model.Club{ID: clubId}
	err = json.NewDecoder(request.Body).Decode(&updatedClub)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		return
	}

	err = database.UpdateClub(updatedClub)
	if err != nil {
		http.Error(writer, "Error updating club", http.StatusInternalServerError)
		return
	}

	response := successResponse{Success: true}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJson)
}

func clubEdit(writer http.ResponseWriter, request *http.Request) error {
	clubId, err := getClubIdFromUrl(writer, request)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return err
	}

	updatedClub := model.Club{ID: clubId}
	err = json.NewDecoder(request.Body).Decode(&updatedClub)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		return err
	}

	err = database.UpdateClub(updatedClub)
	if err != nil {
		http.Error(writer, "Error updating club", http.StatusInternalServerError)
		return err
	}

	response := successResponse{Success: true}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		return err
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJson)

	return nil
}

func ClubView(writer http.ResponseWriter, request *http.Request) {
	user, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	
	clubList := database.GetUserClubList(user.Id)
	club := getClubFromUrl(writer, request)
	
	fmt.Fprintln(writer, "Selected Club:", club.ID, club.Name, club.Description, club.OwnerId, club.DateCreated)
	
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
	user, _ := GetLoggedInUser(request)
	club := getClubFromUrl(writer, request)

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
	user, _ := GetLoggedInUser(request)
	club := getClubFromUrl(writer, request)

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