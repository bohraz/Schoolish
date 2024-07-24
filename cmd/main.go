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

	// Registers a handler to serve files from the /static/ directory, necessary for html files to load css and js files
	staticDir := http.Dir("static")
	staticFS := http.FileServer(staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))

	// Registers predefined function handlers for url requests
	http.HandleFunc("/login/", handlers.ServeFileHandler("static/html/login.html"))
	http.HandleFunc("/signup/", handlers.ServeFileHandler("static/html/register.html"))

	http.HandleFunc("/clubs/", handlers.ClubView)
	http.HandleFunc("/clubs/create", handlers.ServeFileHandler("static/html/createClub.html"))
	http.HandleFunc("/clubs/search/", handlers.ClubSearch)
	http.HandleFunc("/clubs/join/", handlers.ClubJoin)
	http.HandleFunc("/clubs/leave/", handlers.ClubLeave)
	http.Handle("/clubs/edit/", handlers.AuthServeFileHandler("static/html/editClub.html"))

	http.HandleFunc("/api/", handlers.Api)

	// Registers predefined function handlers for url requests that require a user to be logged in
	http.Handle("/secret/", handlers.AuthMiddleware(http.HandlerFunc(handlers.SecretHandler)))

	// Including 127.0.0.1 before port :80 prevents from os requesting permission before every run
	err := http.ListenAndServe("127.0.0.1:80", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
