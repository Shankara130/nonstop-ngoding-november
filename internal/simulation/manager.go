package simulation

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Manager struct {
	clients   map[*websocket.Conn]bool
	broadcast chan interface{}
}

func NewManager() *Manager {
	return &Manager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan interface{}),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case msg := <-m.broadcast:
			for client := range m.clients {
				data, _ := json.Marshal(msg)
				client.WriteMessage(websocket.TextMessage, data)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(manager *Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket error:", err)
		return
	}
	manager.clients[conn] = true

	go func() {
		for {
			msg := map[string]interface{}{
				"time": time.Now().Format("15:04:05"),
				"data": "running",
			}
			manager.broadcast <- msg
			time.Sleep(time.Second)
		}
	}()
}
