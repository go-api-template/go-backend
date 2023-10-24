package common_routes

import (
	"github.com/gin-gonic/gin"
	api "github.com/go-api-template/go-backend/modules/utils/api"
)

type PingRoutesController struct{}

func NewPingRoutesController() PingRoutesController {
	return PingRoutesController{}
}

func (r *PingRoutesController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", func(ctx *gin.Context) {
		api.Ctx(ctx).Ok().SendRaw(gin.H{
			"message": "pong",
		})
	})
}
