package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brunobotter/chat-websocket/dto"
	"github.com/brunobotter/chat-websocket/websocket"
)

func (cw *ClientWrapper) SubscribeAllRooms(ctx context.Context, hub *websocket.Hub) {
	cw.Logger.Info("Iniciando subscriber genérico Redis para todas as salas")
	pubsub := cw.Client.PSubscribe(ctx, "chat:*")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			cw.Logger.Info("Subscriber Redis cancelado pelo contexto")
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			var message dto.Message
			err := json.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				continue
			}

			hub.Broadcast <- message

		}
	}
}

func (cw *ClientWrapper) SaveMessage(ctx context.Context, roomID string, msg dto.Message, maxMessages int) error {
	key := "chat:" + roomID

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// LPUSH adiciona no início da lista
	if err := cw.Client.LPush(ctx, key, payload).Err(); err != nil {
		return err
	}

	// LTRIM mantém apenas as últimas `maxMessages`
	if err := cw.Client.LTrim(ctx, key, 0, int64(maxMessages-1)).Err(); err != nil {
		return err
	}

	if err := cw.Client.Expire(ctx, key, 6*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (cw *ClientWrapper) GetMessages(ctx context.Context, roomID string, limit int) ([]dto.Message, error) {
	key := "chat:" + roomID

	vals, err := cw.Client.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]dto.Message, 0, len(vals))
	for i := len(vals) - 1; i >= 0; i-- {
		var msg dto.Message
		if err := json.Unmarshal([]byte(vals[i]), &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// SaveUnread adiciona uma mensagem privada à lista de mensagens não lidas do usuário
func (cw *ClientWrapper) SaveUnread(ctx context.Context, user string, msg dto.Message) error {
	key := fmt.Sprintf("unread:%s", user)

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := cw.Client.LPush(ctx, key, payload).Err(); err != nil {
		return err
	}

	if err := cw.Client.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

// GetUnreadMessages retorna todas as mensagens não lidas do usuário
func (cw *ClientWrapper) GetUnreadMessages(ctx context.Context, user string) ([]dto.Message, error) {
	key := fmt.Sprintf("unread:%s", user)

	vals, err := cw.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]dto.Message, 0, len(vals))
	for i := len(vals) - 1; i >= 0; i-- { // envia na ordem correta
		var msg dto.Message
		if err := json.Unmarshal([]byte(vals[i]), &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// ClearUnread remove todas as mensagens não lidas do usuário
func (cw *ClientWrapper) ClearUnread(ctx context.Context, user string) error {
	key := fmt.Sprintf("unread:%s", user)
	if err := cw.Client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
