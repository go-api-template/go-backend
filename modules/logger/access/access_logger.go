package access_logger

import (
	"github.com/go-api-template/go-backend/modules/config"
	console_logger "github.com/go-api-template/go-backend/modules/logger/console"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	AccessLogger zerolog.Logger
)

func init() {
	// Get the configuration
	c := config.Config

	// If the debug mode is enabled, we use the console logger
	// Otherwise, we use the file logger
	if c.App.Debug {
		AccessLogger = console_logger.NewConsoleLogger()
	} else {
		// access writer
		fileWriter := &lumberjack.Logger{
			Filename:   c.Logs.Access.Filename,
			MaxSize:    c.Logs.Access.MaxSize,
			MaxBackups: c.Logs.Access.MaxBackups,
			MaxAge:     c.Logs.Access.MaxAge,
			Compress:   c.Logs.Access.Compress,
		}

		// access logger
		AccessLogger = zerolog.New(fileWriter).
			Level(zerolog.TraceLevel).
			With().Timestamp().
			Logger()
	}
}
