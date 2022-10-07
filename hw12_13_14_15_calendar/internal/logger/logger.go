package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logConfig = zap.Config{
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
)

func NewLogger(env, level string) (logger *zap.Logger, err error) {
	switch env {
	case "prod":
		logConfig = zap.NewProductionConfig()
	case "dev":
	default:
		logConfig = zap.NewDevelopmentConfig()
	}

	logLevel := zlevel(level)
	logConfig.Level.SetLevel(logLevel)
	logger, err = logConfig.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return logger, nil
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
