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
	UpdateMe(ctx *gin.Context)
	DeleteMe(ctx *gin.Context)
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
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	httputil.Error
//	@Failure		500		{object}	httputil.Error
//	@Router			/users/me [get]
func (c *UserControllerImpl) GetMe(ctx *gin.Context) {
	// Get the user from the context
	cu, err := middlewares.GetUserFromContext(ctx)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Get the user from the database
	user, err := c.userService.FindById(cu.ID)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}

	// Send the response
	httputil.Ctx(ctx).Ok().Response(user.Response())
}

// UpdateMe godoc
//
//	@Summary		Update information about the connected user
//	@Description	Update information about the connected user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User information"
//	@Success		200		{object}	models.User
//	@Failure		204		{object}	httputil.Message
//	@Failure		400		{object}	httputil.Error
//	@Failure		500		{object}	httputil.Error
//	@Router			/user/me [patch]
func (c *UserControllerImpl) UpdateMe(ctx *gin.Context) {
	// Get the user from the context
	cu, err := middlewares.GetUserFromContext(ctx)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Retrieve the user from the request body
	var payload models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Get the user from the database
	user, err := c.userService.FindById(cu.ID)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Fields to update
	if payload.Name != "" {
		user.Name = payload.Name
	}
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}

	user, err = c.userService.Update(cu.ID, user)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}
	if user == nil {
		httputil.Ctx(ctx).NoContent().Message("unknown user")
		return
	}

	httputil.Ctx(ctx).Ok().Response(user)
}

func (c *UserControllerImpl) DeleteMe(ctx *gin.Context) {
	// Get the user from the context
	cu, err := middlewares.GetUserFromContext(ctx)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Delete the user from the database
	err = c.userService.Delete(cu.ID)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}

	httputil.Ctx(ctx).NoContent().Message("user deleted")
}
