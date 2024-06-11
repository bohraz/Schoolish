package handlers

import (
	"fmt"
	"net/http"
	"os"
	"root/internal/database"
)

// This handler serves the static html login file
func LoginForm(w http.ResponseWriter, r *http.Request) {
	path := "C:/Users/jaibr/OneDrive - University of Missouri/Summer Semester 1/Project/static/html/login.html"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File not found!!!", http.StatusNotFound)
		fmt.Println(err)
		return
	}

	http.ServeFile(w, r, path)
}

// This is a temporary handler that gets the form data, sends it to check login info, and Fprints the result
func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Print("Form failed to parse!", err)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Proof of hashing capabilities, will fully implement when Sign Up page is complete
	//hashedPassword, _ := auth.HashPassword(password)
	//fmt.Println("Hashed password: ", hashedPassword)

	success := database.CheckLoginInfo(username, password)


	// If username and Password succeeded, then echo user exists else redirect back to login
	if success {
		fmt.Fprint(w, "That user exists!")
	} else {
		fmt.Fprintln(w, `<div id="message">That user doesn't exist!</div>`)
		fmt.Fprintln(w, `<script>setTimeout(function() { window.location.href = "/login/"; }, 2000);</script>`)
	}
}