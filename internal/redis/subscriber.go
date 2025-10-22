package redis

import (
	"context"
	"encoding/json"

	"github.com/brunobotter/chat-websocket/internal/dto"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"go.uber.org/zap"
)

func (cw *ClientWrapper) SubscribeMessages(ctx context.Context, roomId string, hub *websocket.Hub) {
	channel := "chat:" + roomId
	cw.Logger.Info("Iniciando subscriber Redis", zap.String("channel", channel))

	pubsub := cw.Client.Subscribe(ctx, channel)
	defer func() {
		_ = pubsub.Close()
		cw.Logger.Info("Subscriber Redis encerrado", zap.String("channel", channel))
	}()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			cw.Logger.Info("Subscriber Redis cancelado pelo contexto", zap.String("channel", channel))
			return
		case msg, ok := <-ch:
			if !ok {
				cw.Logger.Warn("Canal Redis fechado", zap.String("channel", channel))
				return
			}

			// Envia a mensagem para todos os clientes conectados
			var message dto.Message
			_ = json.Unmarshal([]byte(msg.Payload), &message)
			hub.Broadcast <- message
			cw.Logger.Debug("Mensagem recebida do Redis e enviada ao Hub",
				zap.String("channel", channel),
				zap.String("payload", msg.Payload),
			)
		}
	}
}

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
