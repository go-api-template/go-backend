package common_routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
)

type StatusRouteController struct {
}

func NewStatusRoutesController() StatusRouteController {
	return StatusRouteController{}
}

func (r *StatusRouteController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/status", func(ctx *gin.Context) {
		httputil.Ctx(ctx).Ok().SendRaw(gin.H{
			"welcome":        fmt.Sprintf("Welcome to to go-api-template/%s", config.Config.App.Name),
			"version":        config.Config.App.Version,
			"environnement ": config.Config.App.Environnement,
		})
	})
}
