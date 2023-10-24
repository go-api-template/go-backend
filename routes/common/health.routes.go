package common_routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	api "github.com/go-api-template/go-backend/modules/utils/api"
)

type HealthCheckRouteController struct {
}

func NewHealthCheckRoutesController() HealthCheckRouteController {
	return HealthCheckRouteController{}
}

func (r *HealthCheckRouteController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/healthcheck", func(ctx *gin.Context) {
		message := fmt.Sprintf("Welcome to go-api-template/%s", config.Config.App.Name)
		api.Ctx(ctx).Ok().SendRaw(gin.H{
			"message": message,
		})
	})
}
