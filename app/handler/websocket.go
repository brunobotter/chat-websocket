package handler

import (
	"github.com/brunobotter/chat-websocket/redis"
	"github.com/brunobotter/chat-websocket/websocket"
	"github.com/labstack/echo/v4"
)

func WebSocketHandler(hub *websocket.Hub, messageStore redis.MessageStore, publisher redis.Publisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.HandleConnections(hub, c.Response().Writer, c.Request(), messageStore, publisher)
		return nil
	}
}
