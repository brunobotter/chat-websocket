package redis

import (
	"context"
	"encoding/json"

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
