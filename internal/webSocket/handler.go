package webSocket

import (
	"encoding/json"
	"log"
	"server/internal/models"
	"server/internal/services"

	"github.com/gorilla/websocket"
)

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func HandleMessage(conn *websocket.Conn, msg Message) {
	switch msg.Event {
	case "register":
		handleRegister(conn, msg.Data)

	case "movement":
		handleMovement(msg.Data)

	//TODO -> test case remove later
	case "ping":
		if err := conn.WriteJSON(Message{Event: "pong"}); err != nil {
			log.Println("Failed to send pong:", err)
		}
	//TODO -> test case remove later
	case "message":
		log.Println("Message received:", string(msg.Data))

	default:
		log.Println("Unknown event:", msg.Event)
	}
}

func handleRegister(conn *websocket.Conn, data json.RawMessage) {
	var registerData struct {
		UserID string `json:"user_id"`
	}
	if err := json.Unmarshal(data, &registerData); err != nil {
		log.Println("Invalid registration data:", err)
		return
	}

	RegisterClient(conn, registerData.UserID)
}

func handleMovement(data json.RawMessage) {
	var pos models.Position
	if err := json.Unmarshal(data, &pos); err != nil {
		log.Println("Invalid position data:", err)
		return
	}

	Mutex.Lock()
	Positions[pos.UserID] = pos
	Mutex.Unlock()

	services.BroadcastPosition(Positions, Clients)
}
