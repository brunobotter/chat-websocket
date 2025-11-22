package handler

import (
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandlerInterface interface {
	HandleWebSocket(c echo.Context) error
}

type webSocketHandler struct {
	cfg *config.Deps
	hub *websocket.Hub
}

func NewWebSocketHandler(cfg *config.Deps, hub *websocket.Hub) WebSocketHandlerInterface {
	return &webSocketHandler{
		cfg: cfg,
		hub: hub,
	}
}

func (h *webSocketHandler) HandleWebSocket(c echo.Context) error {
	websocket.HandleConnections(h.hub, c.Response().Writer, c.Request(), h.cfg.Redis)
	return nil
}
