package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brunobotter/chat-websocket/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Mapping *Mapping
	Logger  *zap.Logger
	Redis   *redis.ClientWrapper
}

func Init() *Config {

	cfg, err := Read()
	if err != nil {
		panic(fmt.Sprintf("Erro ao ler configuração: %v", err))
	}
	/*redisCfg := redis.RedisConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	}*/

	/*redisClient, err := redis.NewClient(redisCfg, cfg.Logger)
	if err != nil {

	}*/

	deps := &Config{
		Mapping: cfg,
		Logger:  nil,
		Redis:   nil,
	}
	return deps
}

func Read() (*Mapping, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.SetTypeByDefaultValue(true)

	err := viper.ReadInConfig()

	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return nil, err
	}

	config := Mapping{}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
