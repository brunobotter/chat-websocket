package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

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
	logger, err = cfg.Build()

	if err != nil {
		panic(err)
	}
}

func L() *zap.Logger {
	return logger
}
