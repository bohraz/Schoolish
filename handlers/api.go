package handlers

import (
	"log"
	"net/http"
	"strings"
)

var ApiHandlers = make(map[string]func(http.ResponseWriter, *http.Request) error)

func Api(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	pathComponents := strings.Split(path, "/")
	urlRequest := pathComponents[2]

	handler := ApiHandlers[urlRequest]
	if handler == nil {
		log.Println("No handler found for request", urlRequest)
		return
	}
	
	handler(writer, request)
}