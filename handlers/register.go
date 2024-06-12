package handlers

import (
	"fmt"
	"net/http"
	"os"
	"root/internal/auth"
	"root/internal/database"
)

// This handler serves the static html register file
func RegisterForm(writer http.ResponseWriter, request *http.Request) {
	path := "static/html/register.html"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(writer, "File not found!!!", http.StatusNotFound)
		fmt.Println("The error was: ", err)
		return
	}

	http.ServeFile(writer, request, path)
}

// This handler gets the form data, hashes the password, and saves the new user to the database
func RegisterSubmit(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Print("Form failed to parse!", err)
	}

	username := request.FormValue("username")
	password := request.FormValue("password")
	email := request.FormValue("email")
	firstName := request.FormValue("firstname")
	lastName := request.FormValue("lastname")

	if found, str := database.UserFound(username, email); found {
		fmt.Fprintf(writer, "A user with that %s already exists!", str)
		return
	}

	password, err = auth.HashPassword(password)
	if err != nil {
		fmt.Println("There was an error hashing the password!", err)
	}

	accepted := database.CreateUser(username, email, password, firstName, lastName)
	if accepted {
		fmt.Printf("The user %v has been registered!", username)
	}
}
