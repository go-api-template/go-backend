package database_logger

import (
	"github.com/go-api-template/go-backend/modules/config"
	console_logger "github.com/go-api-template/go-backend/modules/logger/console"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DatabaseLogger zerolog.Logger
)

func init() {
	// Get the configuration
	c := config.Config

	// If the debug mode is enabled, we use the console logger
	// Otherwise, we use the file logger
	if c.App.Debug {
		DatabaseLogger = console_logger.NewConsoleLogger()
	} else {
		// access writer
		fileWriter := &lumberjack.Logger{
			Filename:   c.Logs.Database.Filename,
			MaxSize:    c.Logs.Database.MaxSize,
			MaxBackups: c.Logs.Database.MaxBackups,
			MaxAge:     c.Logs.Database.MaxAge,
			Compress:   c.Logs.Database.Compress,
		}

		// access logger
		DatabaseLogger = zerolog.New(fileWriter).
			Level(zerolog.TraceLevel).
			With().Timestamp().
			Logger()
	}
}
