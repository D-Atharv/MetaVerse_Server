package webSocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func HandleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v\n", err)
        return
    }
    defer RemoveClient(conn)

    log.Println("New WebSocket connection established!")

    for {
        var msg Message

        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("Error reading JSON message: %v\n", err)
            log.Println("Closing connection due to invalid message")
            break
        }

        log.Printf("Structured message received: Event=%s, Data=%s", msg.Event, string(msg.Data))
        HandleMessage(conn, msg)

        response := map[string]string{"status": "success", "message": "PONG"}
        if err := conn.WriteJSON(response); err != nil {
            log.Printf("Error writing JSON response: %v\n", err)
            break
        }
    }
}
