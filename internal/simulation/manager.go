package simulation

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Manager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case c := <-m.register:
			m.mu.Lock()
			m.clients[c] = true
			m.mu.Unlock()
			log.Println("Websocket client registered, total:", len(m.clients))
		case c := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[c]; ok {
				delete(m.clients, c)
				close(c.send)
				c.conn.Close()
			}
			m.mu.Unlock()
			log.Println("Websocket client unregistered, total:", len(m.clients))
		case msg := <-m.broadcast:
			m.mu.Lock()
			for c := range m.clients {
				select {
				case c.send <- msg:
				default:
					delete(m.clients, c)
					close(c.send)
					c.conn.Close()
				}
			}
			m.mu.Unlock()
		}
	}
}

func (m *Manager) Broadcast(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("Broadcast marshal error", err)
		return
	}
	select {
	case m.broadcast <- data:
	default:
		log.Println("Broadcast channel full: dropping message")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(manager *Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	manager.register <- client

	go client.writePump()
	go client.readPump(manager)
}

func (c *Client) readPump(m *Manager) {
	defer func() {
		m.unregister <- c
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("readPump read error", err)
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
