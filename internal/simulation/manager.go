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
	for msg := range m.broadcast {
		data, _ := json.Marshal(msg)
		for client := range m.clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("Error send:", err)
				client.Close()
				delete(m.clients, client)
			}
		}
	}
}

func (m *Manager) Broadcast(msg interface{}) {
	m.broadcast <- msg
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
