package handlers

import (
	"fmt"
	"net/http"
	"root/internal/auth"
)

// Temporary static page to test session data persistence across pages
func SecretHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := auth.SESSION_STORE.Get(request, "auth-session")
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
	}

	fmt.Fprintf(writer, "Welcome, %v!", session.Values["userId"])
}