package providers

import (
	"context"

	"github.com/brunobotter/chat-websocket/logger"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/redis"
	"github.com/brunobotter/chat-websocket/websocket"
)

type HubServiceProvider struct{}

func NewHubServiceProvider() *HubServiceProvider {
	return &HubServiceProvider{}
}

func (p *HubServiceProvider) Register(c container.Container) {
	c.Singleton(func(logger logger.Logger, redisClient *redis.ClientWrapper) (*websocket.Hub, error) {
		hub := websocket.NewHub(logger, redisClient)
		go redisClient.SubscribeAllRooms(context.Background(), hub)
		return hub, nil
	})
}
