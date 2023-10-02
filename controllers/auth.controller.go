package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
	"github.com/google/uuid"
	"strings"
)

// IAuthController is the controller for authentification
// It declares the methods that the controller must implement
type IAuthController interface {
	UserSignUp(ctx *gin.Context)
	UserSignIn(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	UserSignOut(ctx *gin.Context)
}

// AuthController is the controller for authentification
// It implements the IAuthController interface
type AuthController struct {
	userService services.IUserService
}

// AuthController implements the IAuthController interface
var _ IAuthController = &AuthController{}

var (
	CtxAccessToken  = "access_token"
	CtxRefreshToken = "refresh_token"
	CtxLoggedIn     = "logged_in"
)

func NewAuthController(userService services.IUserService) IAuthController {
	return &AuthController{userService: userService}
}

// UserSignUp godoc
//
//	@Summary		Sign up a new user
//	@Description	Sign up a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserSignUp	true	"User sign up"
//	@Success		201		{object}	models.UserResponse
//	@Failure		400		{object}	httputil.Error
//	@Failure		409		{object}	httputil.Error
//	@Failure		412		{object}	httputil.Error
//	@Failure		502		{object}	httputil.Error
//	@Router			/auth/signup [post]
func (c *AuthController) UserSignUp(ctx *gin.Context) {
	var payload *models.UserSignUp

	// Bind the request body to the payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Check if the passwords match
	if payload.Password != payload.PasswordConfirm {
		httputil.Ctx(ctx).PreconditionFailed().ErrorMessage("passwords do not match")
		return
	}

	// Sign up the user
	user, err := c.userService.Create(payload)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			httputil.Ctx(ctx).Conflict().Error(err)
			return
		}
		httputil.Ctx(ctx).BadGateway().Error(err)
		return
	}

	// Send the response
	httputil.Ctx(ctx).Created().Response(user.ToResponse())
}

// UserSignIn godoc
//
//	@Summary		Sign in a user
//	@Description	Sign in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			account	body		models.UserSignIn	true	"User credential"
//	@Success		200		{object}	models.AccessToken
//	@Failure		400		{object}	httputil.Error
//	@Failure		404		{object}	httputil.Error
//	@Router			/auth/signin [post]
func (c *AuthController) UserSignIn(ctx *gin.Context) {
	var payload *models.UserSignIn

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Find the user by email
	user, err := c.userService.FindByEmail(payload.Email)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		httputil.Ctx(ctx).NotFound().ErrorMessage("Invalid email or password")
		return
	}

	// Check if the password is correct
	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		httputil.Ctx(ctx).BadRequest().ErrorMessage("Invalid email or Password")
		return
	}

	// todo: check if the user is verified

	// Generate access tokens
	accessToken, err := utils.CreateToken(config.Config.Tokens.Access.ExpiresIn, user.ID, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Generate refresh tokens
	refreshToken, err := utils.CreateToken(config.Config.Tokens.Refresh.ExpiresIn, user.ID, config.Config.Tokens.Refresh.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	ctx.SetCookie(CtxAccessToken, accessToken, config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, refreshToken, config.Config.Tokens.Refresh.MaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, false)

	httputil.Ctx(ctx).Ok().Response(&models.AccessToken{
		AccessToken: accessToken,
		ExpiresIn:   config.Config.Tokens.Access.MaxAge * 60,
		TokenType:   "Bearer",
	})
}

// RefreshAccessToken godoc
//
//	@Summary		Refresh the access token
//	@Description	Refresh the access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.AccessToken
//	@Failure		400	{object}	httputil.Error
//	@Failure		403	{object}	httputil.Error
//	@Failure		404	{object}	httputil.Error
//	@Router			/auth/refresh [get]
func (c *AuthController) RefreshAccessToken(ctx *gin.Context) {

	// Get the refresh token from the cookie
	token, err := ctx.Cookie(CtxRefreshToken)
	if err != nil {
		httputil.Ctx(ctx).Forbidden().ErrorMessage("could not refresh access token")
		return
	}

	// Validate the refresh token
	subscriber, err := utils.ValidateToken(token, config.Config.Tokens.Refresh.PublicKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Convert the subscriber to a UUID
	subscriberUUID, err := uuid.Parse(fmt.Sprint(subscriber))
	if err != nil {
		httputil.Ctx(ctx).BadRequest().ErrorMessage("cannot get user id from token")
		return
	}

	// Get the user from the database
	user, err := c.userService.FindById(subscriberUUID)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		httputil.Ctx(ctx).NotFound().ErrorMessage("the user belonging to this token no longer exists")
		return
	}

	// Generate a new access token
	accessToken, err := utils.CreateToken(config.Config.Tokens.Access.ExpiresIn, user.ID, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	ctx.SetCookie(CtxAccessToken, accessToken, config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, false)

	httputil.Ctx(ctx).Ok().Response(&models.AccessToken{
		AccessToken: accessToken,
	})
}

// UserSignOut godoc
//
//	@Summary		Sign out current user
//	@Description	Sign out current user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	httputil.Message
//	@Router			/auth/signout [get]
func (c *AuthController) UserSignOut(ctx *gin.Context) {
	ctx.SetCookie(CtxAccessToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "", -1, "/", "localhost", false, false)

	httputil.Ctx(ctx).Ok().Success()
}
