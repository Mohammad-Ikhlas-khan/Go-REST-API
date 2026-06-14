package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init initialises the global Zap logger (call once at startup).
func Init(env string) {
	once.Do(func() {
		var cfg zap.Config
		if env == "production" {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}

		var err error
		log, err = cfg.Build()
		if err != nil {
			log = zap.NewNop()
		}
	})
}

// Get returns the singleton logger, initialising with development defaults if needed.
func Get() *zap.Logger {
	if log == nil {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}
		Init(env)
	}
	return log
}

// Sync flushes any buffered log entries. Call at program exit.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
