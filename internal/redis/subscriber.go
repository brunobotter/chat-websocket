package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brunobotter/chat-websocket/internal/dto"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"go.uber.org/zap"
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
				cw.Logger.Warn("Canal Redis fechado")
				return
			}

			var message dto.Message
			err := json.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				cw.Logger.Error("Erro ao desserializar mensagem", zap.Error(err))
				continue
			}

			hub.Broadcast <- message
			cw.Logger.Debug("Mensagem recebida do Redis e enviada ao Hub",
				zap.String("channel", msg.Channel),
				zap.String("payload", msg.Payload),
			)
		}
	}
}

func (cw *ClientWrapper) SaveMessage(ctx context.Context, roomID string, msg dto.Message, maxMessages int) error {
	key := "chat:" + roomID

	payload, err := json.Marshal(msg)
	if err != nil {
		cw.Logger.Error("Falha ao serializar mensagem", zap.Error(err))
		return err
	}

	err = cw.Client.LPush(ctx, key, payload).Err()
	if err != nil {
		cw.Logger.Error("Erro ao salvar mensagem no Redis", zap.String("room", roomID), zap.Error(err))
		return err
	}

	err = cw.Client.LTrim(ctx, key, 0, int64(maxMessages-1)).Err()
	if err != nil {
		cw.Logger.Error("Erro ao limitar mensagens no Redis", zap.String("room", roomID), zap.Error(err))
		return err
	}

	err = cw.Client.Expire(ctx, key, 6*time.Hour).Err()
	if err != nil {
		cw.Logger.Error("Erro ao inserir expiracao no redis", zap.String("room", roomID), zap.Error(err))
		return err
	}

	return nil
}

func (cw *ClientWrapper) GetMessages(ctx context.Context, roomID string, limit int) ([]dto.Message, error) {
	key := "chat:" + roomID

	vals, err := cw.Client.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		cw.Logger.Error("Erro ao buscar mensagens do Redis", zap.String("room", roomID), zap.Error(err))
		return nil, err
	}

	n := len(vals)
	tmp := make([]dto.Message, 0, n)

// Unmarshal in order and reverse at the end for efficiency and clarity
for i := 0; i < n; i++ {
    var msg dto.Message
    if err := json.Unmarshal([]byte(vals[i]), &msg); err != nil {
        cw.Logger.Warn("Falha ao desserializar mensagem", zap.Error(err))
        continue
    }
    tmp = append(tmp, msg)
}
// Reverse slice to return messages in correct order (oldest to newest)
messages := make([]dto.Message, 0, len(tmp))
for i := len(tmp) - 1; i >= 0; i-- {
    messages = append(messages, tmp[i])
}
return messages, nil
}

// SaveUnread adiciona uma mensagem privada à lista de mensagens não lidas do usuário
func (cw *ClientWrapper) SaveUnread(ctx context.Context, user string, msg dto.Message) error {
	key := fmt.Sprintf("unread:%s", user)

	payload, err := json.Marshal(msg)
	if err != nil {
		cw.Logger.Error("Falha ao serializar mensagem não lida", zap.Error(err))
		return err
	}

	err = cw.Client.LPush(ctx, key, payload).Err()
	if err != nil {
		cw.Logger.Error("Erro ao salvar mensagem não lida no Redis", zap.String("user", user), zap.Error(err))
		return err
	}

	err = cw.Client.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		cw.Logger.Error("Erro ao inserir expiracao para mensagens não lidas", zap.String("user", user), zap.Error(err))
		return err
	}

	return nil
}

// GetUnreadMessages retorna todas as mensagens não lidas do usuário
func (cw *ClientWrapper) GetUnreadMessages(ctx context.Context, user string) ([]dto.Message, error) {
	key := fmt.Sprintf("unread:%s", user)

	vals, err := cw.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		cw.Logger.Error("Erro ao buscar mensagens não lidas do Redis", zap.String("user", user), zap.Error(err))
		return nil, err
	}

	n := len(vals)
	tmp := make([]dto.Message, 0, n)
for i := 0; i < n; i++ {
    var msg dto.Message
    if err := json.Unmarshal([]byte(vals[i]), &msg); err != nil {
        cw.Logger.Warn("Falha ao desserializar mensagem não lida", zap.Error(err))
        continue
    }
    tmp = append(tmp, msg)
}
messages := make([]dto.Message, 0, len(tmp))
for i := len(tmp) - 1; i >= 0; i-- {
    messages = append(messages, tmp[i])
}
return messages, nil
}

// ClearUnread remove todas as mensagens não lidas do usuário
func (cw *ClientWrapper) ClearUnread(ctx context.Context, user string) error {
	key := fmt.Sprintf("unread:%s", user)
	err := cw.Client.Del(ctx, key).Err()
	if err != nil {
		cw.Logger.Error("Erro ao limpar mensagens não lidas do Redis", zap.String("user", user), zap.Error(err))
		return err
	}
	return nil
}
