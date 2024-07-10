package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
)

type successResponse struct {
	Success bool `json:"success"`
	Id 	int  `json:"id"`
}

func GetLoggedInUser(request *http.Request) (model.User, error) {
	session, err := auth.SESSION_STORE.Get(request, "auth-session")
	if err != nil {
		msg := fmt.Sprintf("error getting session: %v", err)
		return model.User{}, errors.New(msg)
	}

	userId, ok := session.Values["userId"]
	if !ok {
		msg := fmt.Sprintf("user not found in session: %v", err)
		return model.User{}, errors.New(msg)
	}

	user, err := database.QueryUser(userId)
	if err != nil {
		msg := fmt.Sprintf("error querying user: %v", err)
		return model.User{}, errors.New(msg)
	}

	return user, nil
}

func ServeFileHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.Error(writer, "File not found!!!", http.StatusNotFound)
			fmt.Println(err)
		}
	
		http.ServeFile(writer, request, path)
	}
}

func removeEmptyStringsFromPath(splitPath *[]string) {
    for i := 0; i < len(*splitPath); {
        if (*splitPath)[i] == "" {
            // Remove the element at index i by appending the elements before i with the elements after i
            *splitPath = append((*splitPath)[:i], (*splitPath)[i+1:]...)
        } else {
            i++
        }
    }
}