package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
)

type StatusController struct {
	statusService services.IStatusService
}

func NewStatusController(statusService services.IStatusService) StatusController {
	return StatusController{statusService}
}

// StatusResponse is the struct for the response of the status endpoint
// Create a response message giving information about the server
// -> Welcome to go-api-template/go-backend
// -> Server: version
// -> Database: Type and version
// -> Redis: Type and version
type StatusResponse struct {
	Welcome  string `json:"welcome"`
	Version  string `json:"version"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

func (c *StatusController) Status(ctx *gin.Context) {

	response := StatusResponse{
		Welcome:  fmt.Sprintf("Welcome to to go-api-template/%s", config.Config.App.Name),
		Version:  config.Config.App.Version,
		Database: c.statusService.GetDbConnectionStatus(),
		Redis:    c.statusService.GetRedisConnectionStatus(),
	}

	httputil.Ctx(ctx).Ok().Response(response)
}
