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
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

// AuthController is the controller for authentification
// It declares the methods that the controller must implement
type AuthController interface {
	UserSignUp(ctx *gin.Context)
	Welcome(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	UserSignIn(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	UserSignOut(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	// todo : delete du compte avec email de validation -> anonymisation des données
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

// UserSignUp godoc
//
//	@Summary		Sign up a new user
//	@Description	Sign up a new user
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
func (c *AuthControllerImpl) UserSignUp(ctx *gin.Context) {
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
		if err := c.mailerService.SendVerificationCode(user); err != nil {
			log.Error().Err(err).Msg("Failed to send verification code")
		}
	}()

	// Send the response
	httputil.Ctx(ctx).Created().Response(user.Response())
}

// Welcome godoc
//
//	@Summary		Send welcome email
//	@Description	This re-sends the welcome email to the user if the user is not verified
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"User email"	Format(email)
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	httputil.Error
//	@Failure		403		{object}	httputil.Error
//	@Failure		404		{object}	httputil.Error
//	@Router			/auth/welcome/{email} [get]
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
	user.VerificationCode = utils.Encode(utils.GenerateRandomString(32))

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Send verification code to user email address in background
	go func() {
		if err := c.mailerService.SendVerificationCode(user); err != nil {
			log.Error().Err(err).Msg("Failed to send verification code")
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
//	@Param			verification_code	path		string	true	"verification code sent by email"
//	@Success		200					{object}	httputil.Message
//	@Failure		400					{object}	httputil.Error
//	@Failure		404					{object}	httputil.Error
//	@Router			/auth/verify/{verification_code} [get]
func (c *AuthControllerImpl) VerifyEmail(ctx *gin.Context) {
	// Get the verification code passed in the url
	verificationCode := ctx.Params.ByName("verification_code")

	// Get the user with the verification code
	user, err := c.userService.FindByVerificationCode(verificationCode)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}
	if user == nil {
		message := "the user belonging to this code no longer exists verificationCode " + verificationCode
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
func (c *AuthControllerImpl) UserSignIn(ctx *gin.Context) {
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

	// Generate access tokens
	accessToken, err := utils.CreateToken(config.Config.Tokens.Access.ExpiresIn, user, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Generate refresh tokens
	refreshToken, err := utils.CreateToken(config.Config.Tokens.Refresh.ExpiresIn, user, config.Config.Tokens.Refresh.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Set the cookies
	ctx.SetCookie(CtxAccessToken, accessToken, config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, refreshToken, config.Config.Tokens.Refresh.MaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, false)

	// Send the response
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
func (c *AuthControllerImpl) RefreshAccessToken(ctx *gin.Context) {

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
	accessToken, err := utils.CreateToken(config.Config.Tokens.Access.ExpiresIn, user, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	// Set the cookie
	ctx.SetCookie(CtxAccessToken, accessToken, config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", config.Config.Tokens.Access.MaxAge*60, "/", "localhost", false, false)

	// Send the response
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
func (c *AuthControllerImpl) UserSignOut(ctx *gin.Context) {
	ctx.SetCookie(CtxAccessToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "", -1, "/", "localhost", false, false)

	httputil.Ctx(ctx).Ok().Success()
}

// ForgotPassword godoc
//
//	@Summary		Send a reset token by email
//	@Description	Send a reset token by email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"User email"	Format(email)
//	@Success		200		{object}	httputil.Message
//	@Failure		400		{object}	httputil.Error
//	@Failure		401		{object}	httputil.Error
//	@Router			/auth/password/forgot/{email} [get]
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
	if user.Verified {
		httputil.Ctx(ctx).Forbidden().ErrorMessage("Account already verified")
		return
	}

	// Generate a reset token
	user.ResetPasswordToken = utils.Encode(utils.GenerateRandomString(32))
	user.ResetPasswordAt = time.Now()

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		httputil.Ctx(ctx).BadRequest().Error(err)
		return
	}

	//// Send reset token to user email address in background
	//go func() {
	//	if err := c.mailService.SendResetToken(user, resetToken); err != nil {
	//		log.Error().Err(err).Msg("Failed to send reset token")
	//	}
	//}()

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
//	@Router			/auth/password/reset/{reset_token} [patch]
func (c *AuthControllerImpl) ResetPassword(ctx *gin.Context) {

	//// Get the reset token from the url
	//resetToken := ctx.Params.ByName("reset_password_token")
	//
	//// Get the password from the body
	//var payload *models.UserPasswordConfirmation
	//
	//if err := ctx.ShouldBindJSON(&payload); err != nil {
	//	httputil.Ctx(ctx).BadRequest().Error(err)
	//	return
	//}
	//
	//// Check if the passwords match
	//if payload.Password != payload.PasswordConfirm {
	//	httputil.Ctx(ctx).PreconditionFailed().ErrorMessage("passwords do not match")
	//	return
	//}
	//
	//// Get the user with the reset token
	//user, err := c.userService.FindByResetToken(verificationCode)
	//if err != nil {
	//	httputil.Ctx(ctx).BadRequest().Error(err)
	//	return
	//}
	//if user == nil {
	//	message := "the user belonging to this code no longer exists verificationCode " + verificationCode
	//	httputil.Ctx(ctx).NotFound().ErrorMessage(message)
	//	return
	//}
	//
	//// Password
	//hashedPassword, _ := utils.HashPassword(userCredential.Password)
	//passwordResetToken := utils.Encode(resetToken)

	//// Get the user
	//user, err := c.userService.FindUserByPasswordResetToken(passwordResetToken)
	//if err != nil {
	//	httputil.Ctx(ctx).BadRequest().Error(err)
	//	return
	//}
	//if user == nil {
	//	message := "Invalid email or password"
	//	httputil.Ctx(ctx).NotFound().ErrorMessage(message)
	//	return
	//}
	//if !user.Verified {
	//	httputil.Ctx(ctx).Unauthorized().ErrorMessage("Account not verified")
	//	return
	//}
	//
	//// Set user password
	//if _, err := c.userService.UpdateResetPassword(user.ID, hashedPassword, "", time.Now()); err != nil {
	//	httputil.Ctx(ctx).Forbidden().Error(err)
	//	return
	//}
	//
	//ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	//ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	//ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	//
	//httputil.Ctx(ctx).Ok().Message("Password updated successfully")
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
//	@Router			/auth/password/change [patch]
func (c *AuthControllerImpl) ChangePassword(ctx *gin.Context) {
	//// Get the user from the context
	//currentUser, ok := ctx.Get(middleware.CtxUser)
	//if !ok || currentUser == nil {
	//	httputil.Ctx(ctx).Unauthorized().Message("The user must be logged in")
	//	ctx.Abort()
	//	return
	//}
	//
	//// Cast to User model
	//if _, ok := currentUser.(*models.User); !ok {
	//	httputil.Ctx(ctx).Unauthorized().Message("The user must be logged in")
	//	ctx.Abort()
	//	return
	//}
	//
	//// Convert to user
	//user := currentUser.(*models.User)
	//
	//// Get user credential from body
	//var userCredential *models.ChangePasswordInput
	//if err := ctx.ShouldBindJSON(&userCredential); err != nil {
	//	httputil.Ctx(ctx).BadRequest().Error(err)
	//	return
	//}
	//
	//// Check passwords
	//if userCredential.Password != userCredential.PasswordConfirmation {
	//	httputil.Ctx(ctx).PreconditionFailed().ErrorMessage("passwords do not match")
	//	return
	//}
	//
	//// Password
	//hashedPassword, _ := utils.HashPassword(userCredential.Password)
	//
	//// Set user password
	//if _, err := c.userService.UpdateResetPassword(user.ID, hashedPassword, "", time.Now()); err != nil {
	//	httputil.Ctx(ctx).Forbidden().Error(err)
	//	return
	//}
	//
	//ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	//ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	//ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	//
	//httputil.Ctx(ctx).Ok().Message("Password updated successfully")
}
