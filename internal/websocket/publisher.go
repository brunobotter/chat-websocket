package websocket

import (
	"context"

	"github.com/brunobotter/chat-websocket/internal/dto"
)

type ChatStore interface {
	PublishMessage(ctx context.Context, channel string, msg dto.Message) error
	SaveMessage(ctx context.Context, roomID string, msg dto.Message, limit int) error
	GetMessages(ctx context.Context, roomID string, limit int) ([]dto.Message, error)
	SaveUnread(ctx context.Context, user string, msg dto.Message) error
	GetUnreadMessages(ctx context.Context, user string) ([]dto.Message, error)
	ClearUnread(ctx context.Context, user string) error
}

// Publisher abstracts message publishing for easier testing and decoupling.
type Publisher interface {
	Publish(ctx context.Context, channel string, msg dto.Message) error
}
