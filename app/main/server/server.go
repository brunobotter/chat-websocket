package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brunobotter/chat-websocket/config"
	"github.com/brunobotter/chat-websocket/logger"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/main/server/router"
	"github.com/brunobotter/chat-websocket/websocket"
	"github.com/labstack/echo/v4"
)

type Server struct {
	container container.Container
	config    *config.Config
	logger    logger.Logger
	echo      *echo.Echo
}

func NewServer(c container.Container) (*Server, error) {
	server := &Server{
		container: c,
		echo:      echo.New(),
	}

	c.Resolve(&server.config)
	c.Resolve(&server.logger)

	server.setup()
	return server, nil
}

func (s *Server) setup() {
	s.echo.HideBanner = true

	var cfg *config.Config
	var hub *websocket.Hub
	s.container.Resolve(&cfg)
	s.container.Resolve(&hub)

	router.RegisterRoutes(s.echo, cfg, hub)

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
		address := fmt.Sprintf(":%d", s.config.Server.Port)
		if err := s.echo.Start(address); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatal(err)
		}
	}()
	s.waitForShutdown(ctx)
}
