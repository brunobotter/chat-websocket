package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brunobotter/chat-websocket/internal/auth"
	"github.com/brunobotter/chat-websocket/internal/config"
	"github.com/brunobotter/chat-websocket/internal/logger"
	"github.com/brunobotter/chat-websocket/internal/websocket"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	cfg := config.Init()
	hub := websocket.NewHub(logger.Logger, cfg.Redis)
	ctx := context.Background()
	go hub.Run()

	go cfg.Redis.SubscribeAllRooms(ctx, hub)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		pass := r.FormValue("password")
		if pass != "1234" {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		rooms := []string{"default", "vip"}
		access, _ := auth.GenerateAccessToken(user, rooms)
		refresh, _ := auth.GenerateRefreshToken(user)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		user, err := auth.ValidateRefreshToken(token)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		rooms := []string{"default", "vip"}
		newAccess, _ := auth.GenerateAccessToken(user, rooms)
		json.NewEncoder(w).Encode(map[string]string{"access_token": newAccess})
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleConnections(hub, w, r, cfg.Redis)
	})

	logger.L().Info("🚀 Servidor iniciado", zap.Int("port", cfg.Cfg.Server.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Cfg.Server.Port), nil)
	if err != nil {
		logger.L().Error("❌ Servidor com problema", zap.Error(err))
	}
	defer logger.Logger.Sync()
	defer cfg.Redis.Close()
}
