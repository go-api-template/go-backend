package main_logger

import (
	"github.com/go-mods/zerolog-quick/console/colored"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func init() {
	// console writer
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    true,
		TimeFormat: "2006-01-02 15:04:05",
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.MessageFieldName,
		},
		FormatExtra: colored.Colorize,
	}

	// console logger
	log.Logger = zerolog.New(consoleWriter).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}
