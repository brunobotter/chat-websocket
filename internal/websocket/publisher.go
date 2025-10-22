package websocket

import (
	"context"

	"github.com/brunobotter/chat-websocket/internal/dto"
)

type Publisher interface {
	PublishMessage(ctx context.Context, channel string, msg dto.Message) error
}
