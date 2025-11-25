package redis

import (
	"context"
	"encoding/json"

	"github.com/brunobotter/chat-websocket/dto"
)

func (cw *ClientWrapper) PublishMessage(ctx context.Context, channel string, msg dto.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := cw.Client.Publish(ctx, channel, payload).Err(); err != nil {

		return err
	}

	return nil
}
