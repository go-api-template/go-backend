package main_logger

import (
	console_logger "github.com/go-api-template/go-backend/modules/logger/console"
	"github.com/rs/zerolog/log"
)

func init() {
	// console logger
	log.Logger = console_logger.NewConsoleLogger()
}
