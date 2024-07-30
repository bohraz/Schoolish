package handlers

import (
	"log"
	"net/http"
	"root/internal/model"

	"github.com/gorilla/websocket"
)

var (
	feedUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by returning true
			return true
		},
	}

	postUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by returning true
			return true
		},
	}

	feedClients = make(map[*websocket.Conn]bool)
	postClients = make(map[int][]*websocket.Conn)
)

func FeedChangedApi(writer http.ResponseWriter, request *http.Request) {
	ws, err := feedUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	feedClients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(feedClients, ws)
			break
		}
	}
}

func PostChangedApi(writer http.ResponseWriter, request *http.Request) {
	ws, err := postUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	postId, err := getStrAndConvToInt(writer, request, "id")
	if err != nil {
		return
	}
	
	postClients[postId] = append(postClients[postId], ws)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			removePostClient(ws, postId, postClients)
			break
		}
	}
}

func HandlePostChanged() {
	for {
		post := <-model.PostsBroadcast
		for client := range feedClients {
			err := client.WriteJSON(post)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(feedClients, client)
			}
		}

		for _, client := range postClients[post.Id] {
			err := client.WriteJSON(post)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				removePostClient(client, post.Id, postClients)
			}
		}
	}
}