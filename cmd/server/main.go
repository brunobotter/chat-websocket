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
	hub := websocket.NewHub(logger.Logger, cfg.Redis)
	ctx := context.Background()
	go hub.Run()

	go cfg.Redis.SubscribeAllRooms(ctx, hub)

	router := router.NewRouter(cfg, hub)

	logger.L().Info("🚀 Servidor iniciado", zap.Int("port", cfg.Cfg.Server.Port))
	err := router.Start(fmt.Sprintf(":%d", cfg.Cfg.Server.Port))
	if err != nil {
		logger.L().Error("❌ Servidor com problema", zap.Error(err))
	}

	defer logger.Logger.Sync()
	defer cfg.Redis.Close()
}
