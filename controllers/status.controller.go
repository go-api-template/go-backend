package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
)

type StatusController struct {
	statusService services.StatusService
}

func NewStatusController(statusService services.StatusService) StatusController {
	return StatusController{statusService}
}

// Status godoc
//
// @Summary     Status of the server
// @Description Return information about the server
// @Tags        info
// @Accept      json
// @Produce     json
// @Success     200 {object} controllers.StatusResponse
// @Router      /status[get]
func (c *StatusController) Status(ctx *gin.Context) {

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

	response := StatusResponse{
		Welcome:  fmt.Sprintf("Welcome to to go-api-template/%s", config.Config.App.Name),
		Version:  fmt.Sprintf("%s", config.Config.App.Version),
		Database: fmt.Sprintf("%s", c.statusService.GetDbConnectionStatus()),
		Redis:    fmt.Sprintf("%s", c.statusService.GetRedisConnectionStatus()),
	}

	httputil.Ctx(ctx).Ok().Response(response)
}
