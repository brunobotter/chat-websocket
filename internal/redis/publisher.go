package redis

import (
	"context"
	"encoding/json"

	"github.com/brunobotter/chat-websocket/internal/dto"
	"go.uber.org/zap"
)

// RedisPublisher define interface para publicação de mensagens no Redis
// Facilita testes unitários e mocking
//go:generate mockgen -destination=../mocks/mock_redis_publisher.go -package=mocks github.com/brunobotter/chat-websocket/internal/redis RedisPublisher
// Interface priorizada para facilitar testes e desacoplamento
// Pode ser expandida para outros métodos se necessário

type RedisPublisher interface {
	PublishMessage(ctx context.Context, channel string, msg dto.Message) error
}

// ClientWrapper implementa RedisPublisher
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
