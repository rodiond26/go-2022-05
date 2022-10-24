package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logz zap.Logger
}

func (l Logger) Info(msg string) {
	l.logz.Info(msg)
}

func (l Logger) Debug(msg string) {
	l.logz.Debug(msg)
}

func (l Logger) Warn(msg string) {
	l.logz.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.logz.Error(msg)
}

func NewLogger(env, level string) (logger *Logger, err error) {
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	switch env {
	case "prod":
		logConfig = zap.NewProductionConfig()
	case "dev":
	default:
		logConfig = zap.NewDevelopmentConfig()
	}

	logLevel := zlevel(level)
	logConfig.Level.SetLevel(logLevel)
	zlogger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}
	defer zlogger.Sync()

	return &Logger{
		logz: *zlogger,
	}, nil
}

func zlevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
