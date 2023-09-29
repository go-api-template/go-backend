package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
)

type StatusRouteController struct {
	statusController controllers.StatusController
}

func NewStatusRouteController(statusController controllers.StatusController) StatusRouteController {
	return StatusRouteController{statusController}
}

func (r *StatusRouteController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/status", r.statusController.Status)
}
