package websocket

import (
	"context"

	"github.com/brunobotter/chat-websocket/internal/dto"
)

// ChatStore defines the interface for chat message storage and retrieval operations.
type ChatStore interface {
	// PublishMessage publishes a message to a specific channel.
	PublishMessage(ctx context.Context, channel string, msg dto.Message) error
	// SaveMessage saves a message to a room with a limit on stored messages.
	SaveMessage(ctx context.Context, roomID string, msg dto.Message, limit int) error
	// GetMessages retrieves messages from a room up to the specified limit.
	GetMessages(ctx context.Context, roomID string, limit int) ([]dto.Message, error)
	// SaveUnread saves an unread message for a user.
	SaveUnread(ctx context.Context, user string, msg dto.Message) error
	// GetUnreadMessages retrieves unread messages for a user.
	GetUnreadMessages(ctx context.Context, user string) ([]dto.Message, error)
	// ClearUnread clears all unread messages for a user.
	ClearUnread(ctx context.Context, user string) error
}
