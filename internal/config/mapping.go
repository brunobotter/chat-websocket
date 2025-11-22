package config

import "time"

// Mapping holds the configuration for server and redis.
type Mapping struct {
	Server ServerConfig `mapstructure:"server"`
	Redis  RedisConfig  `mapstructure:"redis"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// RedisConfig holds redis-related configuration.
type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}
