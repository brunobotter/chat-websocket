package handler

import (
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"github.com/labstack/echo/v4"
)

func WebSocketHandler(cfg *config.Deps, hub *websocket.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.HandleConnections(hub, c.Response().Writer, c.Request(), cfg.Redis)
		return nil
	}
}
