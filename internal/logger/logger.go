package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerInterface define os métodos essenciais para facilitar testes e mocks.
type LoggerInterface interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sugar() *zap.SugaredLogger
}

var logger LoggerInterface

func Init() {
	var err error
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	l, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	logger = l
}

// L retorna a instância do logger como interface para facilitar testes.
func L() LoggerInterface {
	return logger
}
