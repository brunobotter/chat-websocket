package handler

import (
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"github.com/labstack/echo/v4"
)

// WebSocketConnector abstracts the connection handling for easier testing and flexibility.
type WebSocketConnector interface {
	HandleConnections(hub *websocket.Hub, w echo.ResponseWriter, r *echo.Request, redis config.RedisClient)
}

// DefaultWebSocketConnector is the production implementation using the actual websocket package.
type DefaultWebSocketConnector struct{}

func (d *DefaultWebSocketConnector) HandleConnections(hub *websocket.Hub, w echo.ResponseWriter, r *echo.Request, redis config.RedisClient) {
	websocket.HandleConnections(hub, w, r, redis)
}

// WebSocketHandler returns an Echo handler using a connector interface for testability.
func WebSocketHandler(cfg *config.Deps, hub *websocket.Hub, connector WebSocketConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		connector.HandleConnections(hub, c.Response().Writer, c.Request(), cfg.Redis)
		return nil
	}
}
