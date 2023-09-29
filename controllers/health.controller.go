package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils/http"
)

type HealthController struct {
}

func NewHealthController() HealthController {
	return HealthController{}
}

// Check godoc
//
// @Summary     Check if the server is online
// @Description Send a welcome message if the server is online
// @Tags        info
// @Accept      json
// @Produce     json
// @Success     200 {object} controllers.StatusResponse
// @Router      /healthcheck [get]
func (c *HealthController) Check(ctx *gin.Context) {
	message := fmt.Sprintf("Welcome to go-api-template/%s", config.Config.App.Name)
	httputil.Ctx(ctx).Ok().Response(message)
}
