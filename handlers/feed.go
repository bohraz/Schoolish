package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
)

// On scroll down get and display next x posts
// If scrolling fast show grey (unloaded) posts, counting how many were scrolled through
// 	On scrolling slow, load visible posts
// 	Then load the posts that were greyed out from closest to current post to furthest (FILO?)

// Step 1 is to develop the above
// Step 2 is to convert it to websocket and have a post feed that updates in real time
// Keep in mind twitter-type feed view meaning only comment count shown from feed view

func CreatePostApi(writer http.ResponseWriter, request *http.Request) {
	var post model.Post
	err := json.NewDecoder(request.Body).Decode(&post)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		return
	}
	
	session, err := auth.SESSION_STORE.Get(request, "auth-session")
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := session.Values["userId"].(int)
	post.UserId = userId

	post.Id, err = database.CreatePost(post)
	if err != nil {
		http.Error(writer, "Error creating post", http.StatusInternalServerError)
		log.Println("Error creating post: ", err)
		return
	} else {
		response := successResponse{Success: true, Id: post.Id}
		responseJson, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, "Error encoding response", http.StatusInternalServerError)
			log.Println("Error encoding response: ", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(responseJson)
	
	}
}

func Feed(writer http.ResponseWriter, request *http.Request) {
	_, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
}