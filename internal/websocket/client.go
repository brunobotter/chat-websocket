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

// ChatStore define interface para persistência e publicação de mensagens
// Facilita mocks e testes unitários
//go:generate mockgen -destination=../../mocks/mock_chatstore.go -package=mocks github.com/brunobotter/chat-websocket/internal/websocket ChatStore
// (adapte o caminho conforme sua estrutura)
type ChatStore interface {
	GetMessages(ctx context.Context, room string, limit int) ([]dto.Message, error)
	PublishMessage(ctx context.Context, channel string, msg dto.Message) error
	SaveMessage(ctx context.Context, room string, msg dto.Message, limit int) error
}

type Client struct {
	Conn   WebsocketConn // Interface para facilitar testes
	Send   chan []byte
	Hub    *Hub
	RoomID string
	User   string
}

// WebsocketConn abstrai métodos do *websocket.Conn para facilitar mocks em testes
//go:generate mockgen -destination=../../mocks/mock_wsconn.go -package=mocks github.com/brunobotter/chat-websocket/internal/websocket WebsocketConn
type WebsocketConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// wsConnWrapper implementa WebsocketConn usando *websocket.Conn real
// Útil para injeção de dependência e testes
// (não exportado)
type wsConnWrapper struct {
	*websocket.Conn
}

func (w *wsConnWrapper) ReadMessage() (int, []byte, error)  { return w.Conn.ReadMessage() }
func (w *wsConnWrapper) WriteMessage(mt int, data []byte) error { return w.Conn.WriteMessage(mt, data) }
func (w *wsConnWrapper) Close() error { return w.Conn.Close() }

func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request, store ChatStore) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

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
		ws.Close()
		return
	}

	client := &Client{
		Conn:   &wsConnWrapper{ws}, // injeta wrapper para interface
		Send:   make(chan []byte, 256),
		Hub:    hub,
		RoomID: room,
		User:   claims.User,
	}
	hub.Register <- client

	if history, err := store.GetMessages(r.Context(), room, 50); err == nil {
		for _, msg := range history {
		// Envia JSON serializado para manter padrão e facilitar parsing no front/testes
			sendMsg, _ := json.Marshal(msg)
			client.Send <- sendMsg
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
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
		// log erro para facilitar debug/teste futuramente (pode ser log.Println ou zap/slog)
		// log.Printf("invalid message: %v", err)
		continue
		}

		timestamp := time.Now().UTC() // padroniza timestamp UTC para facilitar testes e consistência

		msg := dto.Message{
			User:      c.User,
			Content:   incoming.Content,
			Timestamp: timestamp,
			RoomID:    c.RoomID,
			Target:    incoming.Target,
		}

		errPub := store.PublishMessage(context.Background(), "chat:"+c.RoomID, msg)
// log erro se necessário para facilitar troubleshooting/testes futuros (pode ser log.Println)
// if errPub != nil { log.Printf("publish error: %v", errPub) }

// Salva histórico (ignora erro por simplicidade; pode ser tratado/logado em produção)
_ = store.SaveMessage(context.Background(), c.RoomID, msg, 50)
}
}

func (c *Client) writePump() {
// defer close do canal Send pode ser considerado para evitar leaks em testes/extensões futuras.
defer c.Conn.Close()	// fecha conexão ao sair do loop.	for msg := range c.Send {		err := c.Conn.WriteMessage(websocket.TextMessage, msg)	if err != nil {	break
}	}	}	
func RefreshHandler(w http.ResponseWriter, r *http.Request) {	refreshToken := r.Header.Get("Authorization")	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {	refreshToken = refreshToken[7:]	}	user, err := auth.ValidateRefreshToken(refreshToken)	if err != nil {	http.Error(w, "invalid refresh token", http.StatusUnauthorized)	return
}	rooms := []string{"default", "vip"}	newAccessToken, _ := auth.GenerateAccessToken(user, rooms)	w.Header().Set("Content-Type", "application/json")	w.Write([]byte(`{"access_token":"` + newAccessToken + `"}`))	}	func LoginHandler(w http.ResponseWriter, r *http.Request) {	user := r.FormValue("user")	password := r.FormValue("password")	if password != "1234" {	http.Error(w, "invalid credentials", http.StatusUnauthorized)	return
}	rooms := []string{"default", "vip"}	accessToken, _ := auth.GenerateAccessToken(user, rooms)	refreshToken, _ := auth.GenerateRefreshToken(user)	w.Header().Set("Content-Type", "application/json")	w.Write([]byte(`{"access_token":"` + accessToken + `","refresh_token":"` + refreshToken + `"}`))	}	