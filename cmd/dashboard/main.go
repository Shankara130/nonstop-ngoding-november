package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/Shankara130/nonstop-ngoding-november/internal/simulation"
)

var tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

func main() {
	world := simulation.NewWorld(100, 100, 25)
	manager := simulation.NewManager()

	go manager.Run()

	go simulation.SimulationLoop(world, 200*time.Millisecond, func(state []map[string]interface{}) {
		payload := map[string]interface{}{
			"type":   "snapshot",
			"tick":   world.Tick,
			"agents": state,
		}
		manager.Broadcast(payload)
	})

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		simulation.HandleWebSocket(manager, w, r)
	})

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Starting server on :3000")
	http.ListenAndServe(":3000", nil)
}
