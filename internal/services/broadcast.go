package services

import (
	"log"
	"server/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

func BroadcastPosition(positions map[string]models.Position, clients map[*websocket.Conn]string) {

	Mutex := sync.Mutex{}

	Mutex.Lock()
	positionList := make([]models.Position, 0, len(positions))

	for _, pos := range positions {
		positionList = append(positionList, pos)
	}
	Mutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(map[string]interface{}{
			"event": "positions",
			"positions":   positionList,
		})

		if err != nil {
			log.Println("Failed to broadcast positions:", err)
			client.Close()

			delete(clients, client)

		}
	}

}
