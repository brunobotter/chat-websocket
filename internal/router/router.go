package router

import (
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/handler"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"github.com/labstack/echo/v4"
)

func NewRouter(cfg *config.Deps, hub *websocket.Hub) *echo.Echo {
	e := echo.New()

	registerPublicRoutes(e)
	registerProtectedRoutes(e, cfg, hub)

	return e
}

func registerPublicRoutes(e *echo.Echo) {
	e.POST("/login", handler.Login)
	e.POST("/refresh", handler.Refresh)
}

func registerProtectedRoutes(e *echo.Echo, cfg *config.Deps, hub *websocket.Hub) {
	e.GET("/ws", handler.WebSocketHandler(cfg, hub))
}
