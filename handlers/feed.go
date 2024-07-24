package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"root/internal/auth"
	"root/internal/database"
	"root/internal/model"
	"strconv"
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

func GetPostsApi(writer http.ResponseWriter, request *http.Request) {
	_, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	amountStr := request.URL.Query().Get("amount")
    if amountStr == "" {
        http.Error(writer, "Missing amount parameter", http.StatusBadRequest)
        return
    }

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(writer, "Invalid amount parameter", http.StatusBadRequest)
		return
	}

	posts, err := database.GetPosts(amount)
	if err != nil {
		http.Error(writer, "Error getting posts", http.StatusInternalServerError)
		log.Println("Error getting posts: ", err)
		return
	}

	postsJson, err := json.Marshal(posts)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response: ", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(postsJson)
}

func GetPostByIdApi(writer http.ResponseWriter, request *http.Request) {
	_, err := GetLoggedInUser(request)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIdStr := request.URL.Query().Get("id")
	if postIdStr == "" {
		http.Error(writer, "Missing id parameter", http.StatusBadRequest)
		return
	}

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(writer, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	post, err := database.GetPost(postId)
	if err != nil {
		http.Error(writer, "Error getting post", http.StatusInternalServerError)
		log.Println("Error getting post: ", err)
		return
	}

	postJson, err := json.Marshal(post)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response: ", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(postJson)
}

func CreateCommentApi(writer http.ResponseWriter, request *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(request.Body).Decode(&comment)
	if err != nil {
		http.Error(writer, "Error decoding request", http.StatusBadRequest)
		log.Println("Error decoding request: ", err)
		return
	}
	
	session, err := auth.SESSION_STORE.Get(request, "auth-session")
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := session.Values["userId"].(int)
	comment.User, err = database.QueryUser(userId)
	if err != nil {
		http.Error(writer, "Error querying user", http.StatusInternalServerError)
		log.Println("Error querying user: ", err)
		return
	}

	comment.Id, err = database.CreateComment(comment)
	if err != nil {
		http.Error(writer, "Error creating comment", http.StatusInternalServerError)
		log.Println("Error creating comment: ", err)
		return
	} else {
		responseJson, err := json.Marshal(comment)
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