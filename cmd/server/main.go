package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func dbHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/mysql?parseTime=true")
	if err != nil {
		fmt.Println("Error opening mysql: ", err)
	}

	query := `
		SELECT email, firstName, lastName
		FROM app.users
		WHERE userId = ?
	`

	var (
		email string
		first string
		last string
	)

	err = db.QueryRow(query, 1).Scan(&email, &first, &last)
	if err != nil {
		fmt.Println("Error querying row ", err)
	}

	fmt.Fprintln(w, "Row: ", email, first, last)
}

func main() {
	dir := http.Dir("../../static/html")

	// Creates a File-Server connection starting at path dir
	fs := http.FileServer(dir)

	// Registers a handler for /page/ url request and strips /page from url for effective file search
	http.Handle("/page/", http.StripPrefix("/page", fs))

	http.HandleFunc("/test/", dbHandler)

	// Including 127.0.0.1 before port :80 prevents from os requesting permission before every run
	err := http.ListenAndServe("127.0.0.1:80", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}