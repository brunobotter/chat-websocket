package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerInterface define os métodos essenciais do logger
// Pode ser expandida conforme necessário
// Exemplo: Info, Error, Debug, etc.
type LoggerInterface interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	With(fields ...zap.Field) LoggerInterface
}

// loggerImpl implementa LoggerInterface usando zap.Logger
// encapsula o *zap.Logger para uso interno
type loggerImpl struct {
	logger *zap.Logger
}

func (l *loggerImpl) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *loggerImpl) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *loggerImpl) With(fields ...zap.Field) LoggerInterface {
	return &loggerImpl{logger: l.logger.With(fields...)}
}

var (
	Logger LoggerInterface
)

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
	zLogger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	Logger = &loggerImpl{logger: zLogger}
}

func L() LoggerInterface {
	return Logger
}
