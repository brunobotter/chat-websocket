package websocket

import (
	"github.com/brunobotter/chat-websocket/internal/dto"
	"go.uber.org/zap"
)

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan dto.Message
	Register   chan *Client
	Unregister chan *Client
	logger     *zap.Logger
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan dto.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if client.RoomID == "" {
				client.RoomID = "default"
			}
			if _, ok := h.Rooms[client.RoomID]; !ok {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
		case client := <-h.Unregister:
			if clients, ok := h.Rooms[client.RoomID]; ok {
				delete(clients, client)
				close(client.Send)
				if len(clients) == 0 {
					delete(h.Rooms, client.RoomID)
				}
			}
		case msg := <-h.Broadcast:
			//mensagens privadas
			if msg.Target != "" {
				for _, clients := range h.Rooms {
					for client := range clients {
						if client.User == msg.Target {
							select {
							case client.Send <- []byte(msg.Content):
							default:
								close(client.Send)
								delete(clients, client)
							}
						}
					}
				}
				continue
			}
			//broadcast por rooms
			if clients, ok := h.Rooms[msg.RoomID]; ok {
				for client := range clients {
					select {
					case client.Send <- []byte(msg.Content):
					default:
						delete(clients, client)
						close(client.Send)
					}
				}
			}
		}
	}
}
