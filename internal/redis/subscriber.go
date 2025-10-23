package redis

import (
	"context"
	"encoding/json"
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

	// LPUSH adiciona no início da lista
	if err := cw.Client.LPush(ctx, key, payload).Err(); err != nil {
		cw.Logger.Error("Erro ao salvar mensagem no Redis", zap.String("room", roomID), zap.Error(err))
		return err
	}

	// LTRIM mantém apenas as últimas `maxMessages`
	if err := cw.Client.LTrim(ctx, key, 0, int64(maxMessages-1)).Err(); err != nil {
		cw.Logger.Error("Erro ao limitar mensagens no Redis", zap.String("room", roomID), zap.Error(err))
		return err
	}

	if err := cw.Client.Expire(ctx, key, 6*time.Hour).Err(); err != nil {
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

	messages := make([]dto.Message, 0, len(vals))
	for i := len(vals) - 1; i >= 0; i-- {
		var msg dto.Message
		if err := json.Unmarshal([]byte(vals[i]), &msg); err != nil {
			cw.Logger.Warn("Falha ao desserializar mensagem", zap.Error(err))
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
