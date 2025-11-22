package config

import (
	"strings"

	"github.com/brunobotter/chat-websocket/internal/logger"
	"github.com/brunobotter/chat-websocket/internal/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type RedisClient interface {
	// Adicione métodos relevantes usados no código, por exemplo:
	Ping() error
}

type Deps struct {
	Cfg    *Mapping
	Logger Logger
	Redis  RedisClient
}

func Init(logger Logger) *Deps {
	logger.Info("Inicializando configuração")

	cfg, err := read(logger)
	if err != nil {
		logger.Error("Erro ao ler config", zap.Error(err))
	}

	redisCfg := redis.RedisConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	}

	redisClient, err := redis.NewClient(redisCfg, logger)
	if err != nil {
		logger.Error("Nao foi possivel conectar no Redis", zap.Error(err))
	}

	deps := &Deps{
		Cfg:    cfg,
		Logger: logger,
		Redis:  redisClient,
	}
	return deps
}

func read(logger Logger) (*Mapping, error) {
	setupConfig()
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Erro ao ler config", zap.Error(err))
		return nil, err
	}

	config := Mapping{}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func setupConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("cmd/server")
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
