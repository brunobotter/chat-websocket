package router

import (
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/handler"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"github.com/labstack/echo/v4"
)

func NewRouter(cfg *config.Deps, hub *websocket.Hub) *echo.Echo {
	e := echo.New()

	// Rotas públicas
	e.POST("/login", handler.Login)
	e.POST("/refresh", handler.Refresh)

	// Rotas protegidas
	e.GET("/ws", handler.WebSocketHandler(cfg, hub))
	return e
}
