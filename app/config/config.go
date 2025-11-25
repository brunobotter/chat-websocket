package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Init() *Config {

	cfg, err := Read()
	if err != nil {
		panic(fmt.Sprintf("Erro ao ler configuração: %v", err))
	}

	return cfg
}

func Read() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.host", "SERVER_HOST")

	v.BindEnv("redis.addr", "REDIS_ADDR")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("redis.db", "REDIS_DB")
	v.BindEnv("redis.pool_size", "REDIS_POOL_SIZE")
	v.BindEnv("redis.min_idle_conns", "REDIS_MIN_IDLE_CONNS")
	v.BindEnv("redis.dial_timeout", "REDIS_DIAL_TIMEOUT")
	v.BindEnv("redis.read_timeout", "REDIS_READ_TIMEOUT")
	v.BindEnv("redis.write_timeout", "REDIS_WRITE_TIMEOUT")

	v.BindEnv("app_name", "APP_NAME")
	v.BindEnv("env", "ENV")

	// Não precisa ler o arquivo se só usa env, mas pode manter:
	err := v.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return nil, err
	}

	config := Config{}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
