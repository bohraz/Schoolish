package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"root/internal/database"
	"root/internal/model"
	"strconv"
	"strings"
)

func init() {
	ApiHandlers["clubCreate"] = clubCreate
}

func getClubIdFromURL(writer http.ResponseWriter, request *http.Request, idIndex int) int {
	clubId, err := strconv.Atoi(strings.Split(request.URL.Path, "/")[idIndex])
	if err != nil {
		http.Error(writer, "There was an error processing the URL", http.StatusBadRequest)
		return 0
	}

	return clubId
}

func getClubFromURL(writer http.ResponseWriter, request *http.Request, idIndex int) (model.Club) {
	clubId := getClubIdFromURL(writer, request, idIndex)

	club := database.GetClub(clubId)
	if club.ID == 0 {
		http.Error(writer, "Invalid club ID", http.StatusBadRequest)
		return club
	}

	return club
}

type clubCreateRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type clubCreateResponse struct {
	ClubID int `json:"clubId"`
	Name string `json:"name"`
}

func clubCreate(writer http.ResponseWriter, request *http.Request) error {
	user, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return err
	}

	var createRequest clubCreateRequest
	err = json.NewDecoder(request.Body).Decode(&createRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return err
	}


	clubId, err := database.CreateClub(createRequest.Name, createRequest.Description, user.Id)
	if err != nil {
		http.Error(writer, "Error creating club", http.StatusInternalServerError)
		return err
	}

	response := clubCreateResponse{
		ClubID: clubId,
		Name: createRequest.Name,
	}

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
	club := getClubFromURL(writer, request, 2)
	
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
	user, _ := GetLoggedInUser(request)
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

func canEditClub(writer http.ResponseWriter, request *http.Request) bool {
	user, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return false
	}

	clubId := getClubIdFromURL(writer, request, 3)
	role := database.GetUserClubRole(user.Id, clubId)
	if role < 2 {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		msg := fmt.Sprintf("User %d does not have permission to edit club %d", user.Id, clubId)
		log.Println(msg, role)
		return false
	}

	return true
}

func ClubEdit(writer http.ResponseWriter, request *http.Request) {
    clubId := getClubIdFromURL(writer, request, 3)

	if !canEditClub(writer, request) {
		return
	}

    formTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Edit Club</title>
</head>
<body>
    <h1>Edit Club</h1>
    <form action="/clubs/edit/submit/{{.ClubId}}" method="post">
        <label for="clubName">Club Name:</label>
        <input type="text" id="clubName" name="clubName" required><br>
        <label for="clubDescription">Description:</label>
        <textarea id="clubDescription" name="clubDescription" required></textarea><br>
        <button type="submit">Update Club</button>
    </form>
</body>
</html>
`

    // Parse and execute the template
    tmpl, err := template.New("editClubForm").Parse(formTemplate)
    if err != nil {
        http.Error(writer, "Error rendering form", http.StatusInternalServerError)
        return
    }

	clubIdStr := strconv.Itoa(clubId)
	err = tmpl.Execute(writer, map[string]interface{}{
		"ClubId": clubIdStr,
	})
    if err != nil {
        http.Error(writer, "Error rendering form", http.StatusInternalServerError)
        return
    }
}

func ClubEditSubmit(writer http.ResponseWriter, request *http.Request) {
	clubId := getClubIdFromURL(writer, request, 4)
	fmt.Println("Club id: ", clubId)

	if !canEditClub(writer, request) {
		return
	}

	updatedClub := model.Club{
		ID: clubId,
		Name: request.FormValue("clubName"),
		Description: request.FormValue("clubDescription"),
	}

	err := database.UpdateClub(updatedClub)
	if err != nil {
		http.Error(writer, "Error updating club", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(writer, `<div id="message">Club %s updated!</div>`, updatedClub.Name)
	fmt.Fprintf(writer, `<script>setTimeout(function() { window.location.href = "/clubs/%d/"; }, 2000);</script>`, clubId)
}