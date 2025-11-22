package main

import (
	"context"
	"fmt"

	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/logger"
	"github.com/brunobotter/chat-websocket/internal/router"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"go.uber.org/zap"
)

// HubInterface define os métodos esperados do Hub
type HubInterface interface {
	Run()
}

// RedisInterface define os métodos esperados do Redis
type RedisInterface interface {
	SubscribeAllRooms(ctx context.Context, hub websocket.HubSubscriber)
	Close() error
}

// RouterInterface define os métodos esperados do Router
type RouterInterface interface {
	Start(addr string) error
}

func main() {
	logger.Init()
	cfg := config.Init()

	var hub HubInterface = websocket.NewHub(logger.Logger, cfg.Redis)
	ctx := context.Background()
	go hub.Run()

	var redis RedisInterface = cfg.Redis
	go redis.SubscribeAllRooms(ctx, hub.(websocket.HubSubscriber))

	var r RouterInterface = router.NewRouter(cfg, hub)

	logger.L().Info("🚀 Servidor iniciado", zap.Int("port", cfg.Cfg.Server.Port))
	err := r.Start(fmt.Sprintf(":%d", cfg.Cfg.Server.Port))
	if err != nil {
		logger.L().Error("❌ Servidor com problema", zap.Error(err))
	}

	defer logger.Logger.Sync()
	defer redis.Close()
}
