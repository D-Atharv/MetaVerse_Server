package webSocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return 
	}

	defer conn.Close()

	log.Println("New connection established")

	for {
		messageType, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println(err, "failed to read message")
			break
		}

		log.Printf("Received message: %s", msg)

		err = conn.WriteMessage(messageType, msg)

		if err != nil {
			log.Println(err, "failed to write message")
			break
		}
	}

}
