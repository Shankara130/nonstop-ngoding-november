package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/Shankara130/nonstop-ngoding-november/internal/simulation"
	"github.com/Shankara130/nonstop-ngoding-november/internal/websocket"
)

var tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

func main() {
	world := simulation.NewWorld(100.0)

	// spawn zones with different types
	for i := 0; i < 20; i++ {
		world.SpawnAgent(simulation.AgentTypeWorker)
	}
	for i := 0; i < 10; i++ {
		world.SpawnAgent(simulation.AgentTypeExplorer)
	}
	for i := 0; i < 5; i++ {
		world.SpawnAgent(simulation.AgentTypeCollector)
	}
	for i := 0; i < 5; i++ {
		world.SpawnAgent(simulation.AgentTypeGuard)
	}

	// websocket manager
	manager := websocket.NewWSManager()
	go manager.Run()

	// run simulation
	go world.Run(50*time.Millisecond, func(snapshot *simulation.WorldSnapshot) {
		manager.Broadcast(map[string]interface{}{
			"type": "snapshot",
			"data": snapshot,
		})
	})

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(manager, w, r)
	})

	http.HandleFunc("/spawn", func(w http.ResponseWriter, r *http.Request) {
		agentType := r.URL.Query().Get("type")
		if agentType == "" {
			agentType = string(simulation.AgentTypeWorker)
		}
		world.SpawnAgent(simulation.AgentType(agentType))
		w.WriteHeader(http.StatusOK)
	})

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Starting server on :3000")
	http.ListenAndServe(":3000", nil)
}
