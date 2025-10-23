package websocket

import (
	"context"

	"github.com/brunobotter/chat-websocket/internal/dto"
)

type ChatStore interface {
	PublishMessage(ctx context.Context, channel string, msg dto.Message) error
	SaveMessage(ctx context.Context, roomID string, msg dto.Message, limit int) error
	GetMessages(ctx context.Context, roomID string, limit int) ([]dto.Message, error)
}
