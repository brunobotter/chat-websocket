package handler

import (
	"github.com/brunobotter/chat-websocket/websocket"
	"github.com/labstack/echo/v4"
)

func WebSocketHandler(hub *websocket.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.HandleConnections(hub, c.Response().Writer, c.Request(), hub.ChatStore)
		return nil
	}
}
