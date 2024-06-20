package handlers

import (
	"errors"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
)

func GetUser(request *http.Request) (model.User, error) {
	session, err := auth.SESSION_STORE.Get(request, "auth-session")
	if err != nil {
		return model.User{}, err
	}

	userId, ok := session.Values["userId"]
	if !ok {
		return model.User{}, errors.New("user not found in session")
	}

	user, err := database.QueryUser(userId)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}