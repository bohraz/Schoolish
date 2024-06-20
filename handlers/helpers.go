package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
)

func GetUser(request *http.Request) (model.User, error) {
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