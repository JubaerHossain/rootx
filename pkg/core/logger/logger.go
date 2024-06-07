package logger

import (
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

// Init initializes the logger
func Init() error {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		err := Logger.Sync()
		if err != nil {
			Logger.Error("Failed to sync logger", zap.Error(err))
		}
	}
}

// Info logs an informational message
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}
