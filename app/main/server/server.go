package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brunobotter/chat-websocket/config"
	"github.com/brunobotter/chat-websocket/main/adapters"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/main/server/router"
	"github.com/labstack/echo/v4"
)

type Server struct {
	container container.Container
	config    *config.ServerConfig
	//inserir logger
	echo *echo.Echo
}

func NewServer(c container.Container) (*Server, error) {
	server := &Server{
		container: c,
		echo:      echo.New(),
	}

	c.Resolve(&server.config)
	//c.Resolve(&server.logger)

	server.setup()
	return server, nil
}

func (s *Server) setup() {
	s.echo.HideBanner = true

	adapterRouter := adapters.NewEchoRouterAdapter(s.echo)
	s.container.Singleton(func() router.Router {
		return adapterRouter.Group("/api/v1", nil)
	})

}

func (s *Server) waitForShutdown(ctx context.Context) {
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(err)
	}
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		address := fmt.Sprintf(":%d", s.config.Port)
		if err := s.echo.Start(address); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatal(err)
		}
	}()
	s.waitForShutdown(ctx)
}
