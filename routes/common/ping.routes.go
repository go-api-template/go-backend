package common_routes

import (
	"github.com/gin-gonic/gin"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
)

type PingRoutesController struct{}

func NewPingRoutesController() PingRoutesController {
	return PingRoutesController{}
}

func (r *PingRoutesController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", func(ctx *gin.Context) {
		httputil.Ctx(ctx).Ok().Response(gin.H{
			"message": "pong",
		})
	})
}
