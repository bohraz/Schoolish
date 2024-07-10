package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
)

func RegisterApi(writer http.ResponseWriter, request *http.Request) {
	var registerInfo model.User
	err := json.NewDecoder(request.Body).Decode(&registerInfo)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		return
	}
	
	if found, str := database.UserFound(registerInfo.Handle, registerInfo.Email); found {
		fmt.Fprintf(writer, "A user with that %s already exists!", str)
		return
	}

	userId, err := database.CreateUser(registerInfo.Handle, registerInfo.Email, registerInfo.Password, registerInfo.FirstName, registerInfo.LastName)
	if err != nil {
		return
	} else {
		session, err := auth.SESSION_STORE.Get(request, "auth-session")
		if err != nil {
			fmt.Println("There was an error getting the session!", err)
			return
		}

		session.Values["userId"] = userId
		session.Values["authenticated"] = true
		err = session.Save(request, writer)
		if err != nil {
			fmt.Fprint(writer, "There was an error: ", err)
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

		fmt.Printf("The user %v has been registered! Their name is %v %v and their email is %v!", registerInfo.Handle, registerInfo.FirstName, registerInfo.LastName, registerInfo.Email)
	}
}
