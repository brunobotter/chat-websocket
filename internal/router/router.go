package router

import (
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/handler"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"github.com/labstack/echo/v4"
)

// Handler defines the interface for HTTP handlers, facilitating unit testing and mocking.
type Handler interface {
	Login(c echo.Context) error
	Refresh(c echo.Context) error
	WebSocketHandler(cfg *config.Deps, hub *websocket.Hub) echo.HandlerFunc
}

func NewRouter(cfg *config.Deps, hub *websocket.Hub, h Handler) *echo.Echo {
	e := echo.New()

	// Rotas públicas
	e.POST("/login", h.Login)
	e.POST("/refresh", h.Refresh)

	// Rotas protegidas
	e.GET("/ws", h.WebSocketHandler(cfg, hub))
	return e
}
