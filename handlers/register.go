package handlers

import (
	"fmt"
	"net/http"
	"os"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
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

	user := model.User{
		Handle: request.FormValue("username"),
		Email: request.FormValue("email"),
		FirstName: request.FormValue("firstname"),
		LastName: request.FormValue("lastname"),
	}
	password := request.FormValue("password")

	if found, str := database.UserFound(user.Handle, user.Email); found {
		fmt.Fprintf(writer, "A user with that %s already exists!", str)
		return
	}

	userId, err := database.CreateUser(user.Handle, user.Email, password, user.FirstName, user.LastName)
	if err != nil {
		return
	} else {
		user.Id = userId

		session, err := auth.SESSION_STORE.Get(request, "auth-session")
		if err != nil {
			fmt.Println("There was an error getting the session!", err)
		}

		session.Values["user"] = user
		session.Values["authenticated"] = true
		err = session.Save(request, writer)
		if err != nil {
			fmt.Fprint(writer, "There was an error: ", err)
		}

		fmt.Printf("The user %v has been registered! Their name is %v %v and their email is %v!", user.Handle, user.FirstName, user.LastName, user.Email)
	}
}
