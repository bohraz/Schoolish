package main

import (
	"fmt"
	"net/http"
)

func main() {
	dir := http.Dir("../../static/html")

	// Creates a File-Server connection starting at path dir
	fs := http.FileServer(dir)

	// Registers a handler for /page/ url request and strips /page from url for effective file search
	http.Handle("/page/", http.StripPrefix("/page", fs))

	// Including 127.0.0.1 before port :80 prevents from os requesting permission before every run
	err := http.ListenAndServe("127.0.0.1:80", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}