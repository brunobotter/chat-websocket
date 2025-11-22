package config

import "time"

// Config is the interface for application configuration
// Provides access to server and redis configuration
// This allows for easier mocking and extension
// in tests or other implementations.
type Config interface {
	Server() ServerConfiger
	Redis() RedisConfiger
}

type Mapping struct {
	ServerCfg ServerConfig `mapstructure:"server"`
	RedisCfg  RedisConfig  `mapstructure:"redis"`
}

func (m *Mapping) Server() ServerConfiger {
	return &m.ServerCfg
}

func (m *Mapping) Redis() RedisConfiger {
	return &m.RedisCfg
}

// ServerConfiger is the interface for server config
type ServerConfiger interface {
	GetPort() int
	GetHost() string
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func (s *ServerConfig) GetPort() int {
	return s.Port
}
func (s *ServerConfig) GetHost() string {
	return s.Host
}

// RedisConfiger is the interface for redis config
type RedisConfiger interface {
	GetAddr() string
	GetPassword() string
	GetDB() int
	GetPoolSize() int
	GetMinIdleConns() int
	GetDialTimeout() time.Duration
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}

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

func (r *RedisConfig) GetAddr() string           { return r.Addr }
func (r *RedisConfig) GetPassword() string       { return r.Password }
func (r *RedisConfig) GetDB() int                { return r.DB }
func (r *RedisConfig) GetPoolSize() int          { return r.PoolSize }
func (r *RedisConfig) GetMinIdleConns() int      { return r.MinIdleConns }
func (r *RedisConfig) GetDialTimeout() time.Duration  { return r.DialTimeout }
func (r *RedisConfig) GetReadTimeout() time.Duration  { return r.ReadTimeout }
func (r *RedisConfig) GetWriteTimeout() time.Duration { return r.WriteTimeout }
