package services

import (
	"log"
	"server/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

var Mutex = sync.Mutex{}

func BroadcastPosition(positions map[string]models.Position, clients map[*websocket.Conn]string) {
    Mutex.Lock()
    positionList := make([]models.Position, 0, len(positions))
    for _, pos := range positions {
        positionList = append(positionList, pos)
    }
    Mutex.Unlock()

    //TODO-> Log for debugging
    log.Printf("Broadcasting positions: %+v", positionList)

    for client := range clients {
        err := client.WriteJSON(map[string]interface{}{
            "event":     "positions",
            "positions": positionList,
        })

        if err != nil {
            log.Printf("Failed to send positions to client: %v", err)
            client.Close()
            delete(clients, client)
        }
    }
}
