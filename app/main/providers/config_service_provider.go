package providers

import (
	"github.com/brunobotter/chat-websocket/config"
	"github.com/brunobotter/chat-websocket/logger"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/redis"
)

type ConfigServiceProvider struct{}

func NewConfigServiceProvider() *ConfigServiceProvider {
	return &ConfigServiceProvider{}
}

func (p *ConfigServiceProvider) Register(c container.Container) {
	c.Singleton(func() *config.Config {
		cfg := config.Init()

		return cfg
	})
	c.Singleton(func(cfg *config.Config) redis.RedisConfig {
		return redis.RedisConfig{
			Addr:         cfg.Redis.Addr,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			DialTimeout:  cfg.Redis.DialTimeout,
			ReadTimeout:  cfg.Redis.ReadTimeout,
			WriteTimeout: cfg.Redis.WriteTimeout,
			PoolSize:     cfg.Redis.PoolSize,
			MinIdleConns: cfg.Redis.MinIdleConns,
		}
	})
	c.Singleton(func(cfg *config.Config) logger.Logger {
		return logger.NewLoggerZap(cfg.AppName)
	})
}
