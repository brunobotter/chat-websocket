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

func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request, publisher Publisher) {
	ws, _ := upgrader.Upgrade(w, r, nil)
	room := r.URL.Query().Get("room")
	user := r.URL.Query().Get("user")

	client := &Client{
		Conn:   ws,
		Send:   make(chan []byte),
		Hub:    hub,
		RoomID: room,
		User:   user,
	}
	hub.Register <- client

	go client.writePump()
	client.readPump(publisher)
}

func (c *Client) readPump(publisher Publisher) {
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

		// Publica no Redis
		ctx := context.Background()
		_ = publisher.PublishMessage(ctx, "chat:"+c.RoomID, msg)
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
