package main

import (
	// Following modules are called implicitly
	// They must be imported at startup in order to initialize them
	_ "github.com/go-api-template/go-backend/modules"

	// These modules are used in main.go
	"github.com/go-api-template/go-backend/modules/config"
)

func main() {
	println("Version: ", config.Config.App.Version)
	println("BuildDate: ", config.Config.App.BuildDate)
}
