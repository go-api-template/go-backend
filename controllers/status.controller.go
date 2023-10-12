package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
)

// StatusController is the controller for status
// It declares the methods that the controller must implement
type StatusController interface {
	Status(ctx *gin.Context)
}

// StatusControllerImpl is the controller for status
// It implements the StatusController interface
type StatusControllerImpl struct {
	statusService services.StatusService
}

// StatusControllerImpl implements the StatusController interface
var _ StatusController = &StatusControllerImpl{}

func NewStatusController(statusService services.StatusService) StatusController {
	return &StatusControllerImpl{statusService}
}

// StatusResponse is the struct for the response of the status endpoint
// Create a response message giving information about the server
// -> Welcome to go-api-template/go-backend
// -> Version: version
// -> Database: Type and version
// -> Redis: Type and version
type StatusResponse struct {
	Welcome  string `json:"welcome"`
	Version  string `json:"version"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

func (c *StatusControllerImpl) Status(ctx *gin.Context) {

	response := StatusResponse{
		Welcome:  fmt.Sprintf("Welcome to to go-api-template/%s", config.Config.App.Name),
		Version:  config.Config.App.Version,
		Database: c.statusService.GetDbConnectionStatus(),
		Redis:    c.statusService.GetRedisConnectionStatus(),
	}

	httputil.Ctx(ctx).Ok().Response(response)
}
