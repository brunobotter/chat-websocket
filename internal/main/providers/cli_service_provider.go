package providers

import (
	"github.com/brunobotter/chat-websocket/internal/main/cli"
	"github.com/brunobotter/chat-websocket/internal/main/container"
	"github.com/spf13/cobra"
)

type CliServiceProvider struct {
	commands []any
}

func NewCliServiceProvider() *CliServiceProvider {
	return &CliServiceProvider{
		commands: []any{
			cli.NewStartServiceCmd,
		},
	}
}

func (p *CliServiceProvider) Register(c container.Container) {
	c.Singleton(func() *cobra.Command {
		return &cobra.Command{
			Use:   "int",
			Short: "Mobi integradora console command line",
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
