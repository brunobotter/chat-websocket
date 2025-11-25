package cli

import (
	"fmt"

	"github.com/brunobotter/chat-websocket/main/app"
	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/brunobotter/chat-websocket/main/server"
	"github.com/spf13/cobra"
)

func NewStartServiceCmd(app *app.Application, container container.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "start:server",
		Short: "Start Http server",
		Run: func(cmd *cobra.Command, args []string) {
			go func() {
				srv, err := server.NewServer(container)
				if err != nil {
					panic(fmt.Errorf("could not initialize server: %v", &err))
				}
				srv.Run(cmd.Context())
			}()
			app.WaitForShutdownSignal()
		},
	}
}
