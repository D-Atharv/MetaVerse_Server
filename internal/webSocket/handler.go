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

	Mutex.Lock()
	if _, exists := Positions[registerData.UserID]; !exists {
		Positions[registerData.UserID] = models.Position{
			UserID: registerData.UserID,
			X:      0,
			Y:      0,
		}
	}
	Mutex.Unlock()

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
	defer Mutex.Unlock()

	prevPos, exists := Positions[pos.UserID]
	if exists && prevPos.X == pos.X && prevPos.Y == pos.Y {
		// Skip if the position hasn't changed
		return
	}
	Positions[pos.UserID] = pos


	//TODO -> kept alert for now. To shift to pop up and to initiate a video call
	alerts := []string{}                  
	reverseAlerts := map[string][]string{} // Notifications for other players affected by this movement

	for userID, otherPos := range Positions {
		if userID == pos.UserID {
			continue
		}

		if services.CheckProximity(pos, otherPos) {
			// Notify the moving player about other users nearby
			alerts = append(alerts, userID)

			// Queue a notification for the other user
			if _, ok := reverseAlerts[userID]; !ok {
				reverseAlerts[userID] = []string{}
			}
			reverseAlerts[userID] = append(reverseAlerts[userID], pos.UserID)
		}
	}

	if len(alerts) > 0 {
		conn.WriteJSON(map[string]interface{}{
			"event":  "proximity_alert",
			"alerts": alerts,
		})
	}

	for affectedUserID, usersNear := range reverseAlerts {
		for client, clientUserID := range Clients {
			if clientUserID == affectedUserID {
				err := client.WriteJSON(map[string]interface{}{
					"event":  "proximity_alert",
					"alerts": usersNear,
				})
				if err != nil {
					log.Printf("Failed to send proximity alert to %s: %v", affectedUserID, err)
				}
				break
			}
		}
	}

	services.BroadcastPosition(Positions, Clients)
}
