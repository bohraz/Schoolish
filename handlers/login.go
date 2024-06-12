package handlers

import (
	"fmt"
	"net/http"
	"os"
	"root/internal/auth"
	"root/internal/database"
)

// This handler serves the static html login file
func LoginForm(writer http.ResponseWriter, request *http.Request) {
	path := "static/html/login.html"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(writer, "File not found!!!", http.StatusNotFound)
		fmt.Println(err)
		return
	}

	http.ServeFile(writer, request, path)
}

// This is a temporary handler that gets the form data, sends it to check login info, and Fprints the result
func LoginSubmit(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Print("Form failed to parse!", err)
	}

	username := request.FormValue("username")
	passwordFromUser := request.FormValue("password")

	passwordFromDb := database.CheckLoginInfo(username)
	success := auth.CheckPasswordHash(passwordFromUser, passwordFromDb)

	if success {
		fmt.Fprint(writer, "That user exists!")
	} else {
		fmt.Fprintln(writer, `<div id="message">That user doesn't exist!</div>`)
		fmt.Fprintln(writer, `<script>setTimeout(function() { window.location.href = "/login/"; }, 2000);</script>`)
	}
}
