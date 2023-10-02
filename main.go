package main

import (
	"github.com/go-api-template/go-backend/cmd"
	"github.com/go-api-template/go-backend/docs"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-mods/convert"
	"os"
)

func main() {

	// Swagger documentation
	docs.SwaggerInfo.Title = "GoVoit API"
	docs.SwaggerInfo.Description = "API server pour GoVoit"
	docs.SwaggerInfo.Version = config.Config.App.Version
	docs.SwaggerInfo.Host = config.Config.Server.Host
	docs.SwaggerInfo.Schemes = []string{config.Config.Server.Scheme}

	// Prevent port 80 and 443 to be displayed in swagger url
	port := convert.ToUintDef(config.Config.Server.Port, 0)
	if port > 0 && port != 80 && port != 443 {
		docs.SwaggerInfo.Host += ":" + config.Config.Server.Port
	}

	// Execute commands
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
