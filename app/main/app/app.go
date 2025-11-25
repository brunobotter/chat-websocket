package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/brunobotter/chat-websocket/main/container"
	"github.com/spf13/cobra"
)

type Application struct {
	wg         *sync.WaitGroup
	context    context.Context
	cancel     context.CancelFunc
	signalChan chan os.Signal
	container  container.Container
	providers  []any
	silent     bool
}

func NewApplication(providers []any) *Application {
	app := &Application{
		wg:         &sync.WaitGroup{},
		signalChan: make(chan os.Signal, 1),
		container:  container.NewContainer(),
		providers:  providers,
		silent:     false,
	}
	app.container.Singleton(func() *Application {
		return app
	})

	return app
}

func (app *Application) setupContext() {
	ctx, cancel := context.WithCancel(context.Background())
	app.context = ctx
	app.cancel = cancel
	app.container.Singleton(func() context.Context {
		return ctx
	})
}

func (app *Application) registerProviders() {
	for _, provider := range app.providers {
		ref := reflect.ValueOf(provider)
		method := ref.MethodByName("Register")
		if !method.IsValid() {
			continue
		}

		start := time.Now()
		args := []reflect.Value{reflect.ValueOf(app.container)}
		method.Call(args)
		app.logProviderState("registered", ref.Elem().Type().Name(), time.Since(start))
	}
}

func (app *Application) bootstrapProvider() {
	for _, provider := range app.providers {
		ref := reflect.ValueOf(provider)
		method := ref.MethodByName("Boot")
		if !method.IsValid() {
			continue
		}
		app.wg.Add(1)
		go func() {
			start := time.Now()
			defer app.wg.Done()
			app.container.Call(method.Interface())
			app.logProviderState("bootstrapped", ref.Elem().Type().Name(), time.Since(start))
		}()
	}
	app.wg.Wait()
}

func (app *Application) shutdownProviders() {
	for i := len(app.providers) - 1; i >= 0; i-- {
		provider := app.providers[i]
		ref := reflect.ValueOf(provider)
		method := ref.MethodByName("Shutdown")
		if !method.IsValid() {
			continue
		}
		start := time.Now()
		app.container.Call(method.Interface())
		app.logProviderState("stopped", ref.Elem().Type().Name(), time.Since(start))
	}
}

func (app *Application) logProviderState(state string, providerName string, executuonTime time.Duration) {
	if !app.silent {
		log.Printf("%s [%s][%s]", providerName+strings.Repeat(" ", 35-len(providerName)), state, executuonTime)
	}
}

func (app *Application) WaitForShutdownSignal() {
	signal.Notify(app.signalChan, os.Interrupt, syscall.SIGTERM)
	<-app.signalChan
}

func (app *Application) executionRootCommand() {
	var rootCmd *cobra.Command
	app.container.Resolve(&rootCmd)
	err := rootCmd.ExecuteContext(app.context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (app *Application) Bootstrap() {
	if len(os.Args[1:]) == 0 {
		app.silent = true
	}
	app.setupContext()
	app.registerProviders()
	app.bootstrapProvider()
	app.executionRootCommand()
	app.shutdown()
}

func (app *Application) shutdown() {
	app.cancel()
	app.wg.Wait()
	app.shutdownProviders()
}
