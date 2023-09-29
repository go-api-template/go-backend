package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
)

func Cors() gin.HandlerFunc {
	// Cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")

	// Allow all origins while in debug mode
	if config.Config.App.Debug {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = true
	} else {
		// Use the origins from the config file if defined
		// Otherwise allow all origins
		if len(config.Config.Cors.Origins) > 0 {
			corsConfig.AllowOrigins = config.Config.Cors.Origins
		} else {
			corsConfig.AllowAllOrigins = true
		}
		corsConfig.AllowCredentials = true
	}

	return cors.New(corsConfig)
}
