package common_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRouteController struct {
}

func NewSwaggerRouteController() SwaggerRouteController {
	return SwaggerRouteController{}
}

func (r *SwaggerRouteController) NewRoutes(rg *gin.RouterGroup) {

	// Update the title of the swagger doc
	config := func(config *ginSwagger.Config) {
		config.Title = docs.SwaggerInfo.Title
	}

	// Get the swagger wrapper which wraps `http.Handler` into `gin.HandlerFunc`.
	wrapper := ginSwagger.WrapHandler(swaggerFiles.Handler, config)

	rg.GET("/*any", wrapper)
}
