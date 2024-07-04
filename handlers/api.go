package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var ApiHandlers = make(map[string]func(http.ResponseWriter, *http.Request) error)

func Api(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	pathComponents := strings.Split(path, "/")
	removeEmptyStringsFromPath(&pathComponents)

	fmt.Println(len(pathComponents))

	var handler func(http.ResponseWriter, *http.Request) error

	if len(pathComponents) == 2 {
		urlRequest := pathComponents[1]

		handler = ApiHandlers[urlRequest]
		if handler == nil {
			log.Println("No handler found for request 3", urlRequest)
			return
		}
	} else if len(pathComponents) == 4 {
		cameled := strings.ToUpper(pathComponents[2][:1]) + pathComponents[2][1:]
		request := pathComponents[1] + cameled

		log.Println(request)

		handler = ApiHandlers[request]
		if handler == nil {
			log.Println("No handler found for request 5", path)
			return
		}
	} else {
		log.Println("Invalid request", path)
		return
	}

	err := handler(writer, request)
	if err != nil {
		log.Println("Error handling request", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
}