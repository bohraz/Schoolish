package main

import (
	"fmt"
	"net/http"

	"root/handlers"
	"root/internal/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	database.Init()
	defer database.Close()

	router := mux.NewRouter()

	// Registers a handler to serve files from the /static/ directory, necessary for html files to load css and js files
	staticDir := http.Dir("static")
	staticFS := http.FileServer(staticDir)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFS))

	// Registers predefined function handlers for url requests
	router.HandleFunc("/login/", handlers.ServeFileHandler("static/html/login.html"))
	router.HandleFunc("/signup/", handlers.ServeFileHandler("static/html/register.html"))

	router.HandleFunc("/club/{id}", handlers.ClubView).Methods("GET")
	router.HandleFunc("/club/create", handlers.ServeFileHandler("static/html/club_create.html"))
	router.HandleFunc("/clubs/search/", handlers.ClubSearch)
	router.HandleFunc("/club/{id}/join/", handlers.ClubJoin)
	router.HandleFunc("/club/{id}/leave/", handlers.ClubLeave)
	router.Handle("/club/{id}/edit/", handlers.AuthServeFileHandler("static/html/club_edit.html"))

	router.Handle("/feed/", handlers.AuthServeFileHandler("static/html/feed.html"))

	router.HandleFunc("/api/club/create/", handlers.CreateClubApi).Methods("POST")
	router.HandleFunc("/api/club/edit/", handlers.EditClubApi).Methods("POST")
	router.HandleFunc("/api/login/", handlers.LoginApi).Methods("POST")
	router.HandleFunc("/api/register/", handlers.RegisterApi).Methods("POST")
	router.HandleFunc("/api/post/create/", handlers.CreatePostApi).Methods("POST")
	router.HandleFunc("/api/posts", handlers.GetPostsApi).Methods("GET")

	// Including 127.0.0.1 before port :80 prevents from os requesting permission before every run
	err := http.ListenAndServe("127.0.0.1:80", router)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
