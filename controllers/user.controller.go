package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/middlewares"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
	"github.com/google/uuid"
)

// UserController is the controller for user
// It declares the methods that the controller must implement
type UserController interface {
	GetMe(ctx *gin.Context)
	UpdateMe(ctx *gin.Context)
	DeleteMe(ctx *gin.Context)

	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
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
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	httputil.Error
//	@Failure		500	{object}	httputil.Error
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

// DeleteMe godoc
//
//	@Summary		Delete the connected user
//	@Description	Delete the connected user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		204	{object}	httputil.Message
//	@Failure		400	{object}	httputil.Error
//	@Failure		500	{object}	httputil.Error
//	@Router			/user/me [delete]
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

// FindAll godoc
//
//	@Summary		Find all users
//	@Description	Find all users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"Page number"
//	@Param			limit	query		int		false	"Number of items per page"
//	@Param			sort	query		string	false	"Sort by field"
//	@Param			order	query		string	false	"Sort order (asc or desc)"
//	@Param			search	query		string	false	"Search string"
//	@Success		200		{object}	[]models.User
//	@Failure		400		{object}	httputil.Error
//	@Failure		500		{object}	httputil.Error
//	@Router			/users [get]
func (c *UserControllerImpl) FindAll(ctx *gin.Context) {
	// Get the query parameters
	queryParams := httputil.GetFilter(ctx)

	// Find the users
	users, err := c.userService.FindAll(queryParams)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}

	// Send the response
	httputil.Ctx(ctx).Ok().Response(users)
}

// FindById godoc
//
//	@Summary		Find a user by id
//	@Description	Find a user by id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User id"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	httputil.Error
//	@Failure		404	{object}	httputil.Message
//	@Failure		500	{object}	httputil.Error
//	@Router			/users/{id} [get]
func (c *UserControllerImpl) FindById(ctx *gin.Context) {
	// Get the user id
	id := ctx.Param("id")
	// Parse the user id to uuid
	uid, err := uuid.Parse(id)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Find the user
	user, err := c.userService.FindById(uid)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}
	if user == nil {
		httputil.Ctx(ctx).NotFound().Message("user not found")
		return
	}

	// Send the response
	httputil.Ctx(ctx).Ok().Response(user)
}

// Update godoc
//
//	@Summary		Update a user
//	@Description	Update a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"User id"
//	@Param			user	body		models.User	true	"User information"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	httputil.Error
//	@Failure		404		{object}	httputil.Message
//	@Failure		500		{object}	httputil.Error
//	@Router			/users/{id} [patch]
func (c *UserControllerImpl) Update(ctx *gin.Context) {
	// Get the user id
	id := ctx.Param("id")
	// Parse the user id to uuid
	uid, err := uuid.Parse(id)
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
	user, err := c.userService.FindById(uid)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}
	if user == nil {
		httputil.Ctx(ctx).NotFound().Message("user not found")
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
	if payload.Role != "" {
		user.Role = payload.Role
	}

	user, err = c.userService.Update(uid, user)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}
	if user == nil {
		httputil.Ctx(ctx).NotFound().Message("user not found")
		return
	}

	httputil.Ctx(ctx).Ok().Response(user)
}

// Delete godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User id"
//	@Success		204	{object}	httputil.Message
//	@Failure		400	{object}	httputil.Error
//	@Failure		404	{object}	httputil.Message
//	@Failure		500	{object}	httputil.Error
//	@Router			/users/{id} [delete]
func (c *UserControllerImpl) Delete(ctx *gin.Context) {
	// Get the user id
	id := ctx.Param("id")
	// Parse the user id to uuid
	uid, err := uuid.Parse(id)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Delete the user from the database
	err = c.userService.Delete(uid)
	if err != nil {
		httputil.Ctx(ctx).InternalServerError().Error(err)
		return
	}

	httputil.Ctx(ctx).NoContent().Message("user deleted")
}
