package config

import (
	"strings"

	"github.com/brunobotter/chat-websocket/internal/logger"
	"github.com/brunobotter/chat-websocket/internal/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config interface {
	GetRedisConfig() redis.RedisConfig
}

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type RedisClient interface {
	// Adicione métodos relevantes que você usa do redis.ClientWrapper
}

type Deps struct {
	Cfg    Config
	Logger Logger
	Redis  RedisClient
}

func Init() *Deps {
	log := logger.L()
	log.Info("Inicializando configuração")

	cfg, err := read()
	if err != nil {
		log.Error("Erro ao ler config", zap.Error(err))
	}

	redisCfg := cfg.GetRedisConfig()

	redisClient, err := redis.NewClient(redisCfg, log)
	if err != nil {
		log.Error("Nao foi possivel conectar no Redis", zap.Error(err))
	}

	deps := &Deps{
		Cfg:    cfg,
		Logger: log,
		Redis:  redisClient,
	}
	return deps
}

type Mapping struct {
	Redis struct {
		Addr         string
		Password     string
		DB           int
		DialTimeout  int
		ReadTimeout  int
		WriteTimeout int
		PoolSize     int
		MinIdleConns int
	}
}

func (m *Mapping) GetRedisConfig() redis.RedisConfig {
	return redis.RedisConfig{
		Addr:         m.Redis.Addr,
		Password:     m.Redis.Password,
		DB:           m.Redis.DB,
		DialTimeout:  m.Redis.DialTimeout,
		ReadTimeout:  m.Redis.ReadTimeout,
		WriteTimeout: m.Redis.WriteTimeout,
		PoolSize:     m.Redis.PoolSize,
		MinIdleConns: m.Redis.MinIdleConns,
	}
}

func read() (*Mapping, error) {
	setupConfig()
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

func setupConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("cmd/server")
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
