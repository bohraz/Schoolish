package handlers

import (
	"fmt"
	"net/http"
	"os"
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
	password := request.FormValue("password")

	// Proof of hashing capabilities, will fully implement when Sign Up page is complete
	//hashedPassword, _ := auth.HashPassword(password)
	//fmt.Println("Hashed password: ", hashedPassword)

	success := database.CheckLoginInfo(username, password)


	// If username and Password succeeded, then echo user exists else redirect back to login
	if success {
		fmt.Fprint(writer, "That user exists!")
	} else {
		fmt.Fprintln(writer, `<div id="message">That user doesn't exist!</div>`)
		fmt.Fprintln(writer, `<script>setTimeout(function() { window.location.href = "/login/"; }, 2000);</script>`)
	}
}