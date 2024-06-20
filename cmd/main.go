package main

import (
	"fmt"
	"net/http"

	"root/handlers"
	"root/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database.Init()
	defer database.Close()

	dir := http.Dir("static/html")

	// Creates a File-Server connection starting at path dir
	fs := http.FileServer(dir)

	// Registers a handler for /page/ url request and strips /page from url for effective file search
	http.Handle("/page/", http.StripPrefix("/page", fs))

	// Registers predefined function handlers for url requests
	http.HandleFunc("/login/", handlers.LoginForm)
	http.HandleFunc("/login/submit/", handlers.LoginSubmit)
	http.HandleFunc("/signup/", handlers.RegisterForm)
	http.HandleFunc("/signup/submit/", handlers.RegisterSubmit)

	http.HandleFunc("/clubs/", handlers.Club)
	http.HandleFunc("/clubs/create", handlers.ClubCreate)
	http.HandleFunc("/clubs/search/", handlers.ClubSearch)
	http.HandleFunc("/clubs/join/", handlers.ClubJoin)

	// Registers predefined function handlers for url requests that require a user to be logged in
	http.Handle("/secret/", handlers.AuthMiddleware(http.HandlerFunc(handlers.SecretHandler)))

	// Including 127.0.0.1 before port :80 prevents from os requesting permission before every run
	err := http.ListenAndServe("127.0.0.1:80", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
