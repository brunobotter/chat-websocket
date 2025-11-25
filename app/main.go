package main

import (
	"github.com/brunobotter/chat-websocket/main/app"
	"github.com/brunobotter/chat-websocket/main/providers"
)

func main() {
	app.NewApplication(providers.List()).Bootstrap()
	/*
		hub := websocket.NewHub(logger.Logger, cfg.Redis)
		go hub.Run()

		go cfg.Redis.SubscribeAllRooms(ctx, hub)

		router := router.NewRouter(cfg, hub)

		defer logger.Logger.Sync()
		defer cfg.Redis.Close()*/
}
