package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/brunobotter/chat-websocket/internal/dto"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	RoomID string
	User   string
}

func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request, store ChatStore) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	room := r.URL.Query().Get("room")
	user := r.URL.Query().Get("user")

	client := &Client{
		Conn:   ws,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		RoomID: room,
		User:   user,
	}
	hub.Register <- client
	if history, err := store.GetMessages(r.Context(), room, 50); err == nil {
		for _, msg := range history {
			client.Send <- []byte(msg.Content)
		}
	}

	go client.writePump()
	client.readPump(store)
}

func (c *Client) readPump(store ChatStore) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		var incoming dto.Incoming
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
			continue
		}

		msg := dto.Message{
			User:      c.User,
			Content:   incoming.Content,
			Timestamp: time.Now(),
			RoomID:    c.RoomID,
			Target:    incoming.Target,
		}

		ctx := context.Background()
		_ = store.PublishMessage(ctx, "chat:"+c.RoomID, msg)
		_ = store.SaveMessage(ctx, c.RoomID, msg, 50)
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
