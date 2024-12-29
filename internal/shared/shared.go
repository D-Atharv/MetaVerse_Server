package shared

import (
	"encoding/json"
	"server/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients   = make(map[*websocket.Conn]string) // Connected clients
	Positions = make(map[string]models.Position)
	Mutex     = sync.Mutex{} // Thread-safety for shared resources
)

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}