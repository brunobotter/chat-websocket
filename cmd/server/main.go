package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/logger"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	cfg := config.Init()
	hub := websocket.NewHub(logger.Logger)
	ctx := context.Background()
	go hub.Run()

	go cfg.Redis.SubscribeAllRooms(ctx, hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleConnections(hub, w, r, cfg.Redis)
	})

	logger.L().Info("🚀 Servidor iniciado", zap.Int("port", cfg.Cfg.Server.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Cfg.Server.Port), nil)
	if err != nil {
		logger.L().Error("❌ Servidor com problema", zap.Error(err))
	}
	defer logger.Logger.Sync()
	defer cfg.Redis.Close()
}
