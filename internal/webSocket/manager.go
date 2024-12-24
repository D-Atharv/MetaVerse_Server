package webSocket

import (
	"log"
	"server/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients   = make(map[*websocket.Conn]string) // Connected clients
	Positions = make(map[string]models.Position)
	Mutex     = sync.Mutex{} // Thread-safety for shared resources
)

func RegisterClient(conn *websocket.Conn, userID string) {
	Mutex.Lock()
	defer Mutex.Unlock()

	Clients[conn] = userID
	log.Printf("User %s registered", userID)
}

func RemoveClient(conn *websocket.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()

	userID := Clients[conn]
	delete(Clients, conn)
	delete(Positions, userID)
	conn.Close()
	log.Printf("User %s disconnected", userID)
}