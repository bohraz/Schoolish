package handlers

import (
	"log"
	"net/http"
	"root/internal/model"

	"github.com/gorilla/websocket"
)

var commentUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by returning true
		return true
	},
}

var clients = make(map[int][]*websocket.Conn)

func CommentChangedApi(writer http.ResponseWriter, request *http.Request) {
	ws, err := commentUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	postId, err := getStrAndConvToInt(writer, request, "id")
	if err != nil {
		return
	}

	clients[postId] = append(clients[postId], ws)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			removePostClient(ws, postId, clients)
			break
		}
	}
}

func HandleCommentChanged() {
	for {
		comment := <-model.CommentBroadcast
		for _, client := range clients[comment.PostId] {
			err := client.WriteJSON(comment)
			if err != nil {
				log.Println(err)
				removePostClient(client, comment.PostId, clients)
			}
		}
	}
}
