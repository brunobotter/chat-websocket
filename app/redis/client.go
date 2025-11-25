package redis

import (
	"context"
	"time"

	"github.com/brunobotter/chat-websocket/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ClientWrapper struct {
	Client *redis.Client
	Logger logger.Logger
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
}

func NewClient(cfg RedisConfig, logger logger.Logger) (*ClientWrapper, error) {
	logger.Info("Inicializando Redis client", zap.String("addr", cfg.Addr))

	opts := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}

	rdb := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	logger.Info("Redis conectado com sucesso", zap.String("addr", cfg.Addr))

	return &ClientWrapper{
		Client: rdb,
		Logger: logger,
	}, nil
}

func (cw *ClientWrapper) Close() error {
	return cw.Client.Close()
}
