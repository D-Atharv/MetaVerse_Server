package webRTC

import (
	"encoding/json"
	"log"
	"server/internal/shared"

	"github.com/gorilla/websocket"
)

func HandleSignaling(conn *websocket.Conn, msg shared.Message) {
	var signalingData struct {
		To   string          `json:"to"`
		From string          `json:"from"`
		Body json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(msg.Data, &signalingData); err != nil {
		log.Printf("Invalid signaling data: %v\n", err)
		return
	}

	signalingData.From = shared.Clients[conn]
	msg.Data, _ = json.Marshal(signalingData)

	shared.Mutex.Lock()
	defer shared.Mutex.Unlock()

	for client, userID := range shared.Clients {
		if userID == signalingData.To {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Failed to forward signaling message to %s: %v\n", signalingData.To, err)
			}
			return
		}
	}

	log.Printf("Target peer not found: %s\n", signalingData.To)
}