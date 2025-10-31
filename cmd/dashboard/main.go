package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/Shankara130/nonstop-ngoding-november/internal/simulation"
)

var tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

func main() {
	manager := simulation.NewManager()
	go manager.Run()

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
