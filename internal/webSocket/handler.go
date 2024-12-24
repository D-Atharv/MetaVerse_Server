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

    // Log for debugging
    log.Printf("Registered user: %s", registerData.UserID)

    // Initialize the user's position if not already set
    Mutex.Lock()
    if _, exists := Positions[registerData.UserID]; !exists {
        Positions[registerData.UserID] = models.Position{
            UserID: registerData.UserID,
            X:      0,
            Y:      0,
        }
    }
    Mutex.Unlock()

    log.Printf("Broadcasting positions after registration of user %s", registerData.UserID)
    services.BroadcastPosition(Positions, Clients)
}

func handleMovement(data json.RawMessage, conn *websocket.Conn) {
    var pos models.Position
    if err := json.Unmarshal(data, &pos); err != nil {
        log.Println("Invalid position data:", err)
        response := map[string]string{"status": "error", "message": "Invalid position data"}
        conn.WriteJSON(response)
        return
    }

    Mutex.Lock()
    Positions[pos.UserID] = pos
    Mutex.Unlock()

    // TODO:Log for debugging
    log.Printf("Updated position for user %s: X=%.2f, Y=%.2f", pos.UserID, pos.X, pos.Y)

    // Broadcast updated positions to all clients
    log.Println("Broadcasting updated positions to all clients")
    services.BroadcastPosition(Positions, Clients)
}

