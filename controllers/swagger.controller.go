package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerController struct {
}

func NewSwaggerController() SwaggerController {
	return SwaggerController{}
}

func (c *SwaggerController) Swagger() gin.HandlerFunc {
	config := func(config *ginSwagger.Config) {
		config.Title = docs.SwaggerInfo.Title
	}
	return ginSwagger.WrapHandler(swaggerFiles.Handler, config)
}
