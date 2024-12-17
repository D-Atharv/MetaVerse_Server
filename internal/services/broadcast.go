package services

import (
	"log"
	"server/internal/models"

	"github.com/gorilla/websocket"
)

func BroadcastPosition(positions map[string]models.Position, clients map[*websocket.Conn]string) {

	positionList := make([]models.Position,0, len(positions))

	for _, pos := range positions {
		positionList = append(positionList, pos)
	}

	for client := range clients {
		err := client.WriteJSON(map[string]interface{}{
			"messageType": "positions",
			"positions": positionList,
		})

		if err != nil {
			log.Println("Failed to broadcast positions:", err)
			client.Close()
			delete(clients, client)			
		}
	}

}