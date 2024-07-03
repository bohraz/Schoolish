package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Success bool `json:"success"`
}

func init() {
	ApiHandlers["login"] = login
}
// Temporary login submission handler gets form data, checks login info, and prints result
// If successful assigns user info to auth-session cookie
func login(writer http.ResponseWriter, request *http.Request) error {
	var loginInfo loginRequest
	err := json.NewDecoder(request.Body).Decode(&loginInfo)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		return err
	}

	userId, passwordFromDb := database.GetLoginInfo(loginInfo.Username)
	success := auth.CheckPasswordHash(loginInfo.Password, passwordFromDb)
	response := loginResponse{}

	if success {
		session, err := auth.SESSION_STORE.Get(request, "auth-session")
		if err != nil {
			log.Println("There was an error getting the session!", err)
		}

		session.Values["userId"] = userId
		session.Values["authenticated"] = true
		session.Save(request, writer)

		response.Success = true
		log.Println("User logged in!")
	} else {
		response.Success = false
		log.Println("User login failed!")
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
