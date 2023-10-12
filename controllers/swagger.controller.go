package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerController is the controller for swagger
// It declares the methods that the controller must implement
type SwaggerController interface {
	Swagger() gin.HandlerFunc
}

// SwaggerControllerImpl is the controller for swagger
// It implements the SwaggerController interface
type SwaggerControllerImpl struct{}

// SwaggerControllerImpl implements the SwaggerController interface
var _ SwaggerController = &SwaggerControllerImpl{}

func NewSwaggerController() SwaggerController {
	return &SwaggerControllerImpl{}
}

func (c *SwaggerControllerImpl) Swagger() gin.HandlerFunc {
	config := func(config *ginSwagger.Config) {
		config.Title = docs.SwaggerInfo.Title
	}
	return ginSwagger.WrapHandler(swaggerFiles.Handler, config)
}
