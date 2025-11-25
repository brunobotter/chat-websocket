package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/brunobotter/chat-websocket/auth"
	"github.com/brunobotter/chat-websocket/dto"
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

	// 1. Pegando token do header Authorization
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		ws.Close()
		return
	}

	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims, err := auth.ValidateAccessToken(tokenStr)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		ws.Close()
		return
	}

	// 2. Sala desejada
	room := r.URL.Query().Get("room")
	if room == "" {
		room = "default"
	}
	// 3. Verifica se usuário tem acesso à sala
	authorized := false
	for _, r := range claims.Rooms {
		if r == room {
			authorized = true
			break
		}
	}
	if !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		ws.Close()
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

	// 5. Envia histórico
	if history, err := store.GetMessages(r.Context(), room, 50); err == nil {
		for _, msg := range history {
			client.Send <- []byte(msg.Content)
		}
	}

	// Mensagem de boas-vindas
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

	// Salas permitidas (geralmente buscaria no DB)
	rooms := []string{"default", "vip"}
	newAccessToken, _ := auth.GenerateAccessToken(user, rooms)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"access_token":"` + newAccessToken + `"}`))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// usuário + senha (aqui só um exemplo simples)
	user := r.FormValue("user")
	password := r.FormValue("password")

	if password != "1234" { // exemplo fixo, no real você verifica DB
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// salas que o usuário pode acessar
	rooms := []string{"default", "vip"}

	accessToken, _ := auth.GenerateAccessToken(user, rooms)
	refreshToken, _ := auth.GenerateRefreshToken(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"access_token":"` + accessToken + `","refresh_token":"` + refreshToken + `"}`))
}
