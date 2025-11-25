package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	SetCommonFields(commonFields map[string]any)
	InfoF(format string, args ...interface{})
	Info(args ...interface{})
	ErrorF(format string, args ...interface{})
	Log(msg string)
	WithContext(ctx context.Context) Logger
	WithFields(fields map[string]any) Logger
}

type loggerZap struct {
	appName      string
	level        string
	logger       *zap.Logger
	commonFields []any
}

func NewLoggerZap(appName string) Logger {
	var config zap.Config
	var zapLogger *zap.Logger

	config.Encoding = "json"
	config.EncoderConfig = buildEncondingConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	zapLogger, _ = config.Build()

	j := &loggerZap{
		appName:      appName,
		level:        config.Level.String(),
		logger:       zapLogger,
		commonFields: []any{},
	}

	j.SetCommonFields(map[string]any{
		"application_name": appName,
	})
	return j
}

func (l *loggerZap) SetCommonFields(commonfields map[string]any) {
	for key, value := range commonfields {
		l.commonFields = append(l.commonFields, zap.Any(key, value))
	}
}

func buildEncondingConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		TimeKey:        "timestamp",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (l *loggerZap) InfoF(format string, args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Infof(format, args...)
}

func (l *loggerZap) Info(args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Info(args...)
}

func (l *loggerZap) ErrorF(format string, args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Errorf(format, args...)
}

func (l *loggerZap) Log(msg string) {
	l.Info(msg)
}

func (l *loggerZap) WithContext(ctx context.Context) Logger {
	fields := map[string]any{}
	return l.WithFields(fields)
}

func (l *loggerZap) WithFields(fields map[string]any) Logger {
	newLogger := &loggerZap{
		appName:      l.appName,
		level:        l.level,
		logger:       l.logger,
		commonFields: append([]any{}, l.commonFields...),
	}

	for key, value := range fields {
		newLogger.commonFields = append(newLogger.commonFields, zap.Any(key, value))
	}

	return newLogger
}
