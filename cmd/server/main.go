package main

import (
	"log"
	"net/http"

	"github.com/brunobotter/chat-websocket/internal/websocket"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleConnections(hub, w, r)
	})

	log.Println("🚀 Servidor iniciado em :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
