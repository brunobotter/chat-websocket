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

func main() {
	logger.Init()
	cfg := config.Init()

	// Use interfaces for hub and redis to facilitate testing and decoupling
	var hub websocket.HubInterface = websocket.NewHub(logger.Logger, cfg.Redis)
	var redisClient config.RedisInterface = cfg.Redis

	ctx := context.Background()
	go hub.Run()

	go redisClient.SubscribeAllRooms(ctx, hub)

	r := router.NewRouter(cfg, hub)

	logger.L().Info("🚀 Servidor iniciado", zap.Int("port", cfg.Cfg.Server.Port))
	err := r.Start(fmt.Sprintf(":%d", cfg.Cfg.Server.Port))
	if err != nil {
		logger.L().Error("❌ Servidor com problema", zap.Error(err))
	}

	defer logger.Logger.Sync()
	defer redisClient.Close()
}
