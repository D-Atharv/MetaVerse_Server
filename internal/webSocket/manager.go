package webSocket

import (
	"log"
	"server/internal/shared"

	"github.com/gorilla/websocket"
)

func RegisterClient(conn *websocket.Conn, userID string) {
	shared.Mutex.Lock()
	defer shared.Mutex.Unlock()

	shared.Clients[conn] = userID
	log.Printf("User %s registered", userID)
}

func RemoveClient(conn *websocket.Conn) {
    shared.Mutex.Lock()
    defer shared.Mutex.Unlock()

    userID := shared.Clients[conn]
    delete(shared.Clients, conn)
    delete(shared.Positions, userID)

    for client := range shared.Clients {
        err := client.WriteJSON(map[string]interface{}{
            "event":   "disconnect",
            "user_id": userID,
        })
        if err != nil {
            log.Printf("Failed to notify client about disconnection: %v", err)
        }
    }
    conn.Close()
    log.Printf("User %s disconnected", userID)
}