package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils/http"
)

// HealthController is the controller for health
// It declares the methods that the controller must implement
type HealthController interface {
	Check(ctx *gin.Context)
}

// HealthControllerImpl is the controller for health check
// It implements the HealthController interface
type HealthControllerImpl struct{}

// HealthControllerImpl implements the HealthController interface
var _ HealthController = &HealthControllerImpl{}

func NewHealthController() HealthController {
	return &HealthControllerImpl{}
}

func (c *HealthControllerImpl) Check(ctx *gin.Context) {
	message := fmt.Sprintf("Welcome to go-api-template/%s", config.Config.App.Name)
	httputil.Ctx(ctx).Ok().Message(message)
}
