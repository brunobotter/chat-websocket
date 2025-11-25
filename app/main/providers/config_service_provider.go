package providers

import (
	"github.com/brunobotter/chat-websocket/config"
	"github.com/brunobotter/chat-websocket/logger"
	"github.com/brunobotter/chat-websocket/main/container"
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
	c.Singleton(func(cfg *config.Config) logger.Logger {
		return logger.NewLoggerZap(cfg.Mapping.AppName)
	})
}
