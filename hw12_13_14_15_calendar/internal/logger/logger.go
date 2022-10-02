package logger

import (
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/config"
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

func NewLogger(mainConfig *config.Config) (logger *zap.Logger, err error) {
	switch mainConfig.Environment.Type {
	case "prod":
		logConfig = zap.NewProductionConfig()
	case "dev":
	default:
		logConfig = zap.NewDevelopmentConfig()
	}

	logLevel := level(mainConfig.Logger.Level)
	logConfig.Level.SetLevel(logLevel)
	logger, err = logConfig.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return logger, nil
}

func level(lvl string) zapcore.Level {
	switch lvl {
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
