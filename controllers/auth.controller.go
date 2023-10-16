package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/middlewares"
	"github.com/go-api-template/go-backend/modules/utils"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

// AuthController is the controller for authentification
// It declares the methods that the controller must implement
type AuthController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	SignOut(ctx *gin.Context)
	Welcome(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	RefreshTokens(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

// AuthControllerImpl is the controller for authentification
// It implements the AuthController interface
type AuthControllerImpl struct {
	userService   services.UserService
	mailerService services.MailerService
}

// AuthControllerImpl implements the AuthController interface
var _ AuthController = &AuthControllerImpl{}

var (
	CtxAccessToken  = "access_token"
	CtxRefreshToken = "refresh_token"
	CtxLoggedIn     = "logged_in"
)

func NewAuthController(userService services.UserService, mailerService services.MailerService) AuthController {
	return &AuthControllerImpl{userService: userService, mailerService: mailerService}
}

// SignUp godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserSignUp	true	"User sign up"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	httputil.Error
//	@Failure		409		{object}	httputil.Error
//	@Failure		412		{object}	httputil.Error
//	@Failure		502		{object}	httputil.Error
//	@Router			/auth/signup [post]
func (c *AuthControllerImpl) SignUp(ctx *gin.Context) {
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

	// Send verification code to user email address in background
	go func() {
		if err := c.mailerService.SendVerificationToken(user); err != nil {
			log.Error().Err(err).Msg("Failed to send verification code")
		}
	}()

	// Send the response
	httputil.Ctx(ctx).Created().Response(user.Response())
}

// SignIn godoc
//
//	@Summary		Sign in a user
//	@Description	Sign in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			account	body		models.UserSignIn	true	"User credential"
//	@Success		201		{object}	models.AccessToken
//	@Failure		400		{object}	httputil.Error
//	@Failure		404		{object}	httputil.Error
//	@Router			/auth/signin [post]
func (c *AuthControllerImpl) SignIn(ctx *gin.Context) {
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
	// Check if the account is verified
	if !user.Verified {
		httputil.Ctx(ctx).BadRequest().ErrorMessage("Account not verified")
		return
	}

	// Sign in the user
	c.signin(ctx, user)
}

func (c *AuthControllerImpl) signin(ctx *gin.Context, user *models.User) {
	now := time.Now().UTC()

	// Generate access tokens
	accessToken, err := utils.CreateToken(now, config.Config.Tokens.Access.MaxAge, user, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	accessTokenExpiresIn := config.Config.Tokens.Access.MaxAge * 60
	accessTokenExpiresAt := now.Add(time.Duration(config.Config.Tokens.Access.MaxAge) * time.Minute)

	// Generate refresh tokens
	refreshToken, err := utils.CreateToken(now, config.Config.Tokens.Refresh.MaxAge, user, config.Config.Tokens.Refresh.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	refreshTokenExpiresIn := config.Config.Tokens.Refresh.MaxAge * 60

	// Set the cookies
	ctx.SetCookie(CtxAccessToken, accessToken, accessTokenExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, refreshToken, refreshTokenExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", accessTokenExpiresIn, "/", "localhost", false, false)

	// Send the response
	httputil.Ctx(ctx).Created().Response(&models.AccessToken{
		AccessToken: accessToken,
		ExpiresIn:   accessTokenExpiresIn,
		ExpiresAt:   accessTokenExpiresAt.Unix(),
		TokenType:   "Bearer",
	})
}

// SignOut godoc
//
//	@Summary		Sign out current user
//	@Description	Sign out current user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	httputil.Message
//	@Failure		401	{object}	httputil.Error
//	@Router			/auth/signout [get]
func (c *AuthControllerImpl) SignOut(ctx *gin.Context) {
	// Get the user from the context
	user, err := middlewares.GetUserFromContext(ctx)
	if err != nil || user == nil {
		httputil.Ctx(ctx).Unauthorized().Message("The user must be logged in")
		ctx.Abort()
		return
	}

	// Sign out the user
	c.signOut(ctx, user)
}

func (c *AuthControllerImpl) signOut(ctx *gin.Context, user *models.User) {
	// Clear the cookies
	ctx.SetCookie(CtxAccessToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "", -1, "/", "localhost", false, false)

	// Send the response
	httputil.Ctx(ctx).Ok().Success()
}

// Welcome godoc
//
//	@Summary		Send welcome email
//	@Description	This re-sends the welcome email to the user if the user is not verified
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"User email"	Format(email)
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	httputil.Error
//	@Failure		403		{object}	httputil.Error
//	@Failure		404		{object}	httputil.Error
//	@Router			/auth/welcome/{email} [post]
func (c *AuthControllerImpl) Welcome(ctx *gin.Context) {

	// Get the user email passed in the url
	email := ctx.Params.ByName("email")

	// Get the user from the database by email
	user, err := c.userService.FindByEmail(email)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		message := "unknown user"
		httputil.Ctx(ctx).NotFound().ErrorMessage(message)
		return
	}
	if user.Verified {
		httputil.Ctx(ctx).Forbidden().ErrorMessage("Account already verified")
		return
	}

	// Generate a new Verification Code
	user.VerificationToken = utils.Encode(utils.GenerateRandomString(32))

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Send verification code to user email address in background
	go func() {
		if err := c.mailerService.SendVerificationToken(user); err != nil {
			log.Error().Err(err).Msg("Failed to send the verification token")
		}
	}()

	// Send the response
	httputil.Ctx(ctx).Created().Response(user.Response())
}

// VerifyEmail godoc
//
//	@Summary		Verify email address
//	@Description	Verify email address from verification code sent by email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"verification code sent by email"
//	@Success		200		{object}	httputil.Message
//	@Failure		400		{object}	httputil.Error
//	@Failure		404		{object}	httputil.Error
//	@Router			/auth/verify/{token} [get]
func (c *AuthControllerImpl) VerifyEmail(ctx *gin.Context) {
	// Get the verification code passed in the url
	verificationToken := ctx.Params.ByName("token")

	// Get the user with the verification code
	user, err := c.userService.FindByVerificationToken(verificationToken)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		message := "the user belonging to this code no longer exists verificationToken " + verificationToken
		httputil.Ctx(ctx).NotFound().ErrorMessage(message)
		return
	}

	// Set the user as verified
	user.Verified = true

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	httputil.Ctx(ctx).Ok().Message("Email verified successfully")
}

// RefreshTokens godoc
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
func (c *AuthControllerImpl) RefreshTokens(ctx *gin.Context) {
	now := time.Now().UTC()

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
	accessToken, err := utils.CreateToken(now, config.Config.Tokens.Access.MaxAge, user, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	accessTokenExpiresIn := config.Config.Tokens.Access.MaxAge * 60
	accessTokenExpiresAt := now.Add(time.Duration(config.Config.Tokens.Access.MaxAge) * time.Minute)

	// Set the cookie
	ctx.SetCookie(CtxAccessToken, accessToken, accessTokenExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", accessTokenExpiresIn, "/", "localhost", false, false)

	// Send the response
	httputil.Ctx(ctx).Ok().Response(&models.AccessToken{
		AccessToken: accessToken,
		ExpiresIn:   accessTokenExpiresIn,
		ExpiresAt:   accessTokenExpiresAt.Unix(),
		TokenType:   "Bearer",
	})
}

// ForgotPassword godoc
//
//	@Summary		Send a reset token by email
//	@Description	Send a reset token by email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"User email"	Format(email)
//	@Success		201		{object}	httputil.Message
//	@Failure		400		{object}	httputil.Error
//	@Failure		401		{object}	httputil.Error
//	@Router			/auth/forgot-password/{email} [post]
func (c *AuthControllerImpl) ForgotPassword(ctx *gin.Context) {

	// Get the user email passed in the url
	email := ctx.Params.ByName("email")

	// Get the user from the database by email
	user, err := c.userService.FindByEmail(email)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		message := "unknown user"
		httputil.Ctx(ctx).NotFound().ErrorMessage(message)
		return
	}
	if !user.Verified {
		httputil.Ctx(ctx).Unauthorized().ErrorMessage("Account not verified")
		return
	}

	// Generate a reset token
	user.SetResetToken(utils.Encode(utils.GenerateRandomString(32)))

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Send verification code to user email address in background
	go func() {
		if err := c.mailerService.SendResetToken(user); err != nil {
			log.Error().Err(err).Msg("Failed to send the reset token")
		}
	}()

	// Send the response
	httputil.Ctx(ctx).Created().Response(user.Response())
}

// ResetPassword godoc
//
//	@Summary		Reset the user password
//	@Description	Reset the user password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			password	body		models.UserPasswordConfirmation	true	"New password"
//	@Success		200			{object}	httputil.Message
//	@Failure		400			{object}	httputil.Error
//	@Failure		401			{object}	httputil.Error
//	@Failure		403			{object}	httputil.Error
//	@Failure		404			{object}	httputil.Error
//	@Failure		412			{object}	httputil.Error
//	@Router			/auth/reset-password/{token} [patch]
func (c *AuthControllerImpl) ResetPassword(ctx *gin.Context) {

	// Get the reset token from the url
	resetToken := ctx.Params.ByName("token")

	// Get the password from the body
	var payload *models.UserPasswordConfirmation
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Check if the passwords match
	if payload.Password != payload.PasswordConfirm {
		httputil.Ctx(ctx).PreconditionFailed().ErrorMessage("passwords do not match")
		return
	}

	// Get the user with the reset token
	user, err := c.userService.FindByResetPasswordToken(resetToken)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		message := "the user belonging to this code no longer exists. The reset token was " + resetToken
		httputil.Ctx(ctx).NotFound().ErrorMessage(message)
		return
	}

	// This is the new password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	user.Password = hashedPassword

	// Clear the reset token
	user.SetResetToken("")

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Sign out the user
	c.signOut(ctx, user)
}

// ChangePassword godoc
//
//	@Summary		Change the user password
//	@Description	Change the user password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			password	body		models.UserPasswordConfirmation	true	"New password"
//	@Success		200			{object}	httputil.Message
//	@Failure		400			{object}	httputil.Error
//	@Failure		401			{object}	httputil.Error
//	@Failure		403			{object}	httputil.Error
//	@Failure		404			{object}	httputil.Error
//	@Failure		412			{object}	httputil.Error
//	@Router			/auth/change-password [post]
func (c *AuthControllerImpl) ChangePassword(ctx *gin.Context) {

	// Get the user from the context
	user, err := middlewares.GetUserFromContext(ctx)
	if err != nil || user == nil {
		httputil.Ctx(ctx).Unauthorized().Message("The user must be logged in")
		ctx.Abort()
		return
	}

	// Get the new password from the body
	var payload *models.UserPasswordConfirmation
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Check if the passwords match
	if payload.Password != payload.PasswordConfirm {
		httputil.Ctx(ctx).PreconditionFailed().ErrorMessage("passwords do not match")
		return
	}

	// This is the new password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	user.Password = hashedPassword

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Sign in the user
	c.signin(ctx, user)
}
