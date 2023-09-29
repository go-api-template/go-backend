package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
)

type HealthRouteController struct {
	healthController controllers.HealthController
}

func NewHealthRouteController(HealthController controllers.HealthController) HealthRouteController {
	return HealthRouteController{HealthController}
}

func (r *HealthRouteController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/healthcheck", r.healthController.Check)
}
