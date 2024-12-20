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
		handleMovement(msg.Data, conn) 

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

func handleMovement(data json.RawMessage, conn *websocket.Conn) {
	var pos models.Position
	if err := json.Unmarshal(data, &pos); err != nil {
		log.Println("Invalid position data:", err)
		response := map[string]string{"status": "error", "message": "Invalid position data"}
		conn.WriteJSON(response)
		return
	}

	// Update the position in the server
	Mutex.Lock()
	Positions[pos.UserID] = pos
	Mutex.Unlock()

	response := map[string]interface{}{
		"status":  "success",
		"message": "Position received",
		"data": map[string]interface{}{
			"x": pos.X,
			"y": pos.Y,
		},
	}
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Failed to respond to client:", err)
	}

	//broadcast updated positions to all clients
	services.BroadcastPosition(Positions, Clients)
}

