package redis

import (
	"context"
	"encoding/json"

	"github.com/brunobotter/chat-websocket/dto"
	"go.uber.org/zap"
)

func (cw *ClientWrapper) PublishMessage(ctx context.Context, channel string, msg dto.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		cw.Logger.Error("Falha ao serializar mensagem", zap.Error(err))
		return err
	}

	if err := cw.Client.Publish(ctx, channel, payload).Err(); err != nil {
		cw.Logger.Error("Erro ao publicar mensagem no Redis",
			zap.String("channel", channel),
			zap.Error(err),
		)
		return err
	}

	cw.Logger.Info("Mensagem publicada no Redis",
		zap.String("channel", channel),
		zap.String("payload", string(payload)),
	)
	return nil
}
