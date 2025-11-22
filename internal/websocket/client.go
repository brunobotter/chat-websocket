package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/brunobotter/chat-websocket/internal/auth"
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
	defer func() {
		if ws != nil {
			ws.Close()
		}
	}()

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims, err := auth.ValidateAccessToken(tokenStr)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	room := r.URL.Query().Get("room")
	if room == "" {
		room = "default"
	}
	authorized := false
	for _, r := range claims.Rooms {
		if r == room {
			authorized = true
			break
		}
	}
	if !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	client := &Client{
		Conn:   ws,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		RoomID: room,
		User:   claims.User,
	}
	hub.Register <- client

	if history, err := store.GetMessages(r.Context(), room, 50); err == nil {
		for _, msg := range history {
			client.Send <- []byte(msg.Content)
		}
	}

	msg, _ := json.Marshal(map[string]string{"msg": "connected to " + room})
	client.Send <- msg

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
		err = json.Unmarshal(msgBytes, &incoming)
		if err != nil {
			continue
		}

		timestamp := time.Now()

		target := incoming.Target

		dtoMsg := dto.Message{
			User:      c.User,
			Content:   incoming.Content,
			Timestamp: timestamp,
			RoomID:    c.RoomID,
			Target:    target,
		}

		errPub := store.PublishMessage(context.Background(), "chat:"+c.RoomID, dtoMsg)

		errSave := store.SaveMessage(context.Background(), c.RoomID, dtoMsg, 50)

// Optionally log errors here if needed.
_ = errPub
_ = errSave

// Optionally: handle errors if needed.
// if errPub != nil || errSave != nil { ... }

// Optionally: add rate limiting or validation here.

// End for loop.	}
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

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
refreshToken := r.Header.Get("Authorization")
if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
refreshToken = refreshToken[7:]
}
user, err := auth.ValidateRefreshToken(refreshToken)
if err != nil {
http.Error(w, "invalid refresh token", http.StatusUnauthorized)
return
}
rooms := []string{"default", "vip"}
newAccessToken, _ := auth.GenerateAccessToken(user, rooms)
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"access_token":"` + newAccessToken + `"}`))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
user := r.FormValue("user")
password := r.FormValue("password")
if password != "1234" {
http.Error(w, "invalid credentials", http.StatusUnauthorized)
return
}
rooms := []string{"default", "vip"}
accessToken, _ := auth.GenerateAccessToken(user, rooms)
refreshToken, _ := auth.GenerateRefreshToken(user)
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"access_token":"` + accessToken + `","refresh_token":"` + refreshToken + `"}`))
}
