package main

import (
	"github.com/brunobotter/chat-websocket/main/app"
	"github.com/brunobotter/chat-websocket/main/providers"
)

func main() {
	app.NewApplication(providers.List()).Bootstrap()
	/*logger.Init()
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
	defer cfg.Redis.Close()*/
}
