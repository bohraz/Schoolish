package handlers

import (
	"net/http"
	"root/internal/auth"
)

// Authentication middleware for pages that require a user to be logged in
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := auth.SESSION_STORE.Get(request, "auth-session")
		if err != nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        	return
		}

		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			http.Error(writer, "Unauthorizied", http.StatusUnauthorized)
			return
		}

		// Call the next hamdler, if authentication is successful
		next.ServeHTTP(writer, request)
	})
}