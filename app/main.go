package main

import (
	"github.com/brunobotter/chat-websocket/main/app"
	"github.com/brunobotter/chat-websocket/main/providers"
)

func main() {
	app.NewApplication(providers.List()).Bootstrap()

}
