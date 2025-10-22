package config

import (
	"os"

	"github.com/brunobotter/chat-websocket/internal/logger"
	"github.com/brunobotter/chat-websocket/internal/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Deps struct {
	Cfg    *Mapping
	Logger *zap.Logger
	Redis  *redis.ClientWrapper
	//Svc contract.ServiceManager
}

func Init() *Deps {
	logger.L().Info("Inicializando configuração")

	profile := os.Getenv("PROFILE")
	cfg, err := read(profile)
	if err != nil {
		logger.L().Error("Erro ao ler config", zap.Error(err))
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

	redisClient, err := redis.NewClient(redisCfg, logger.L())
	if err != nil {
		logger.L().Error("Nao foi possivel conectar no Redis", zap.Error(err))

	}

	deps := &Deps{
		Cfg:    cfg,
		Logger: logger.L(),
		Redis:  redisClient,
	}
	return deps
}

func read(profile string) (*Mapping, error) {
	setupConfig(profile)
	err := viper.ReadInConfig()
	if err != nil {
		logger.L().Error("Erro ao ler config", zap.Error(err))
		return nil, err
	}

	config := Mapping{}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func setupConfig(profile string) {
	viper.AutomaticEnv()
	viper.SetConfigName(profile)
	viper.AddConfigPath(".")
}
