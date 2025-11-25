package adapters

import (
	"github.com/brunobotter/chat-websocket/main/server/router"
	"github.com/labstack/echo/v4"
)

type echoRouterAdapter struct {
	echo *echo.Group
}

func NewEchoRouterAdapter(e *echo.Echo) router.Router {
	return &echoRouterAdapter{
		echo: e.Group(""),
	}
}

func (a *echoRouterAdapter) Group(prefix string, group func(group router.RouteGroup)) router.RouteGroup {
	echoGroup := a.echo.Group(prefix)
	groupRouter := &echoRouterAdapter{
		echo: echoGroup,
	}

	if groupRouter != nil {
		group(groupRouter)
	}
	return groupRouter

}
