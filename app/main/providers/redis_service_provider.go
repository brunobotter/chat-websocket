package providers

import (
	"github.com/brunobotter/chat-websocket/logger"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/redis"
)

type RedisServiceProvider struct{}

func NewRedisServiceProvider() *RedisServiceProvider {
	return &RedisServiceProvider{}
}

func (p *RedisServiceProvider) Register(c container.Container) {
	c.Singleton(func(redisConfig redis.RedisConfig, logger logger.Logger) (*redis.ClientWrapper, error) {
		return redis.NewClient(redisConfig, logger)
	})
}
