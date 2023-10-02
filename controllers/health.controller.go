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

func (c *HealthController) Check(ctx *gin.Context) {
	message := fmt.Sprintf("Welcome to go-api-template/%s", config.Config.App.Name)
	httputil.Ctx(ctx).Ok().Message(message)
}
