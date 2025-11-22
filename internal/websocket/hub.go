package websocket

import (
	"context"
	"encoding/json"

	"github.com/brunobotter/chat-websocket/internal/dto"
	"go.uber.org/zap"
)

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan dto.Message
	Register   chan *Client
	Unregister chan *Client
	logger     *zap.Logger
	ChatStore  ChatStore
}

func NewHub(logger *zap.Logger, chatStore ChatStore) *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan dto.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		logger:     logger,
		ChatStore:  chatStore,
	}
}

func (h *Hub) Run() {
	ctx := context.Background()
	for {
		select {
		case client := <-h.Register:
			h.handleRegister(ctx, client)
		case client := <-h.Unregister:
			h.handleUnregister(client)
		case msg := <-h.Broadcast:
			h.handleBroadcast(ctx, msg)
		}
	}
}

func (h *Hub) handleRegister(ctx context.Context, client *Client) {
	if client.RoomID == "" {
		client.RoomID = "default"
	}
	if _, ok := h.Rooms[client.RoomID]; !ok {
		h.Rooms[client.RoomID] = make(map[*Client]bool)
	}
	h.Rooms[client.RoomID][client] = true

	if h.ChatStore != nil {
		go func(c *Client) {
			unread, err := h.ChatStore.GetUnreadMessages(ctx, c.User)
			if err != nil {
				h.logger.Error("Falha ao buscar mensagens não lidas", zap.String("user", c.User), zap.Error(err))
				return
			}
			for _, msg := range unread {
				payload, _ := json.Marshal(msg)
				c.Send <- payload
			}
			_ = h.ChatStore.ClearUnread(ctx, c.User)
		}(client)
	}
}

func (h *Hub) handleUnregister(client *Client) {
	if clients, ok := h.Rooms[client.RoomID]; ok {
		delete(clients, client)
		close(client.Send)
		if len(clients) == 0 {
			delete(h.Rooms, client.RoomID)
		}
	}
}

func (h *Hub) handleBroadcast(ctx context.Context, msg dto.Message) {
	if msg.Target != "" {
		h.sendPrivateMessage(ctx, msg)
	} else {
		h.broadcastToRoom(ctx, msg)
	}
}

func (h *Hub) sendPrivateMessage(ctx context.Context, msg dto.Message) {
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
// Salva mensagem como não lida no Redis
if h.ChatStore != nil {
	h.ChatStore.SaveUnread(ctx, msg.Target, msg)
}
}

func (h *Hub) broadcastToRoom(ctx context.Context, msg dto.Message) {
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
// Salva a mensagem da sala no Redis
if h.ChatStore != nil {
h.ChatStore.SaveMessage(ctx, msg.RoomID, msg, 50)
}
}
