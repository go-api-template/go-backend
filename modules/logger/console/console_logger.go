package console_logger

import (
	"github.com/go-mods/zerolog-quick/console/colored"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"os"
)

var isTerm = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

func NewConsoleLogger() zerolog.Logger {

	// console writer to be used in the console logger
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    true,
		TimeFormat: "2006-01-02 15:04:05",
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
	}

	// If the output is a terminal, colorize the output
	if isTerm {
		consoleWriter.NoColor = false
		consoleWriter.FormatExtra = colored.Colorize
	}

	l := zerolog.New(consoleWriter).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	return l
}
