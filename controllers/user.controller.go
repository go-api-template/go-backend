package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/middlewares"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
)

// UserController is the controller for user
// It declares the methods that the controller must implement
type UserController interface {
	GetMe(ctx *gin.Context)
}

// UserControllerImpl is the controller for user
// It implements the UserController interface
type UserControllerImpl struct {
	userService services.UserService
}

// UserControllerImpl implements the UserController interface
var _ UserController = &UserControllerImpl{}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{userService: userService}
}

// GetMe godoc
//
//	@Summary		Get the current user
//	@Description	Get the current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Router			/users/me [get]
func (u *UserControllerImpl) GetMe(ctx *gin.Context) {
	user := ctx.MustGet(middlewares.CtxUser).(*models.User)
	if user == nil {
		httputil.Ctx(ctx).BadRequest().ErrorMessage("cannot get user")
		return
	}

	// Send the response
	httputil.Ctx(ctx).Ok().Response(user.Response())
}
