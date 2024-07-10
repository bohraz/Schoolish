package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
)

func LoginApi(writer http.ResponseWriter, request *http.Request)  {
	var loginInfo model.User
	err := json.NewDecoder(request.Body).Decode(&loginInfo)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		return
	}

	userId, passwordFromDb := database.GetLoginInfo(loginInfo.Handle)
	success := auth.CheckPasswordHash(loginInfo.Password, passwordFromDb)
	response := successResponse{}

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
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJson)
}
