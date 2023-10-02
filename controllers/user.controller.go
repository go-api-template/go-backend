package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/middlewares"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
)

// IUserController is the controller for user
// It declares the methods that the controller must implement
type IUserController interface {
	GetMe(ctx *gin.Context)
}

// UserController is the controller for user
// It implements the IUserController interface
type UserController struct {
	userService services.IUserService
}

// UserController implements the IUserController interface
var _ IUserController = &UserController{}

func NewUserController(userService services.IUserService) IUserController {
	return &UserController{userService: userService}
}

// GetMe godoc
//
//	@Summary		Get the current user
//	@Description	Get the current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Router			/users/me [get]
func (u *UserController) GetMe(ctx *gin.Context) {
	user := ctx.MustGet(middlewares.CtxUser).(*models.User)
	if user == nil {
		httputil.Ctx(ctx).BadRequest().ErrorMessage("cannot get user")
		return
	}

	httputil.Ctx(ctx).Ok().Response(user.ToResponse())
}
