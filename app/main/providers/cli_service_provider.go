package providers

import (
	"fmt"

	"github.com/brunobotter/chat-websocket/main/app"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/main/server"
	"github.com/brunobotter/chat-websocket/websocket"
	"github.com/spf13/cobra"
)

type CliServiceProvider struct {
	commands []any
}

func NewCliServiceProvider() *CliServiceProvider {
	return &CliServiceProvider{}
}

func (p *CliServiceProvider) Register(c container.Container) {
	var hub *websocket.Hub
	c.Resolve(&hub)
	go hub.Run()
	c.Singleton(func(app *app.Application, container container.Container) *cobra.Command {
		return &cobra.Command{
			Use:   "int",
			Short: "Mobi integradora console command line",
			Run: func(cmd *cobra.Command, args []string) {
				srv, err := server.NewServer(container)
				if err != nil {
					panic(fmt.Errorf("could not initialize server: %v", &err))
				}
				srv.Run(cmd.Context())
				app.WaitForShutdownSignal()
			},
		}
	})
}

func (p *CliServiceProvider) Boot(c container.Container, root *cobra.Command) {
	p.registerCommands(c, root)
}

func (p *CliServiceProvider) registerCommands(c container.Container, root *cobra.Command) {
	for _, constructor := range p.commands {
		cmd, ok := c.Call(constructor).(*cobra.Command)
		if !ok {
			panic("command must be a instance of cobra")
		}
		root.AddCommand(cmd)
	}
}
