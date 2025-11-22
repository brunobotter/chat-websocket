package config

import "time"

// ServerConfigInterface define métodos para facilitar testes e mocks
// Facilita a substituição de implementações em testes unitários
// Pode ser expandida conforme necessidade
// Exemplo: GetPort(), GetHost()
type ServerConfigInterface interface {
	GetPort() int
	GetHost() string
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func (s ServerConfig) GetPort() int  { return s.Port }
func (s ServerConfig) GetHost() string { return s.Host }

// RedisConfigInterface para facilitar testes unitários e mocks
// Pode ser expandida conforme necessidade
// Exemplo: GetAddr(), GetPassword(), etc.
type RedisConfigInterface interface {
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

func (r RedisConfig) GetAddr() string         { return r.Addr }
func (r RedisConfig) GetPassword() string     { return r.Password }
func (r RedisConfig) GetDB() int              { return r.DB }
func (r RedisConfig) GetPoolSize() int        { return r.PoolSize }
func (r RedisConfig) GetMinIdleConns() int    { return r.MinIdleConns }
func (r RedisConfig) GetDialTimeout() time.Duration  { return r.DialTimeout }
func (r RedisConfig) GetReadTimeout() time.Duration  { return r.ReadTimeout }
func (r RedisConfig) GetWriteTimeout() time.Duration { return r.WriteTimeout }

type Mapping struct {
	Server ServerConfig `mapstructure:"server"`
	Redis  RedisConfig  `mapstructure:"redis"`
}
