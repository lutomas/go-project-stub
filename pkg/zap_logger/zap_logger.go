package zap_logger

import (
	"sync"

	"github.com/lutomas/go-project-stub/pkg/config"
	"go.uber.org/zap"
)

var logger *zap.Logger
var once sync.Once

func GetInstance() *zap.Logger {
	once.Do(func() {
		logger = createLogger(&config.Logger{
			DisableStackTrace:    false,
			ProductionFormatting: true,
		})
	})
	return logger
}

func GetInstanceFromConfig(cfg *config.Logger) *zap.Logger {
	once.Do(func() {
		logger = createLogger(cfg)
	})
	return logger
}

func createLogger(cfg *config.Logger) *zap.Logger {
	var zapLogger *zap.Logger
	var err error
	if cfg.ProductionFormatting {
		zapLogger, err = zap.NewProduction()
	} else {
		loggingCfg := zap.NewDevelopmentConfig()
		loggingCfg.DisableStacktrace = cfg.DisableStackTrace
		zapLogger, err = loggingCfg.Build()
	}
	if err != nil {
		panic(err)
	}

	return zapLogger
}

func GetCLIInstance() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Development = true
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
	// zapLogger.DI
}
