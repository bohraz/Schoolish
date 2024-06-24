package handlers

import (
	"fmt"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
)

// Temporary login submission handler gets form data, checks login info, and prints result
// If successful assigns user info to auth-session cookie
func LoginSubmit(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Print("Form failed to parse!", err)
	}

	username := request.FormValue("username")
	passwordFromUser := request.FormValue("password")

	userId, passwordFromDb := database.GetLoginInfo(username)
	success := auth.CheckPasswordHash(passwordFromUser, passwordFromDb)

	if success {
		session, err := auth.SESSION_STORE.Get(request, "auth-session")
		if err != nil {
			fmt.Println("There was an error getting the session!", err)
		}

		session.Values["userId"] = userId
		session.Values["authenticated"] = true
		session.Save(request, writer)

		fmt.Println("User logged in!")
	} else { 
		fmt.Fprintln(writer, `<div id="message">That user doesn't exist!</div>`)
		fmt.Fprintln(writer, `<script>setTimeout(function() { window.location.href = "/login/"; }, 2000);</script>`)
	}
}
