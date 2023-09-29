package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
)

type SwaggerRouteController struct {
	SwaggerController controllers.SwaggerController
}

func NewSwaggerRouteController(SwaggerController controllers.SwaggerController) SwaggerRouteController {
	return SwaggerRouteController{SwaggerController}
}

func (r *SwaggerRouteController) NewRoutes(rg *gin.RouterGroup) {
	rg.GET("/*any", r.SwaggerController.Swagger())
}
