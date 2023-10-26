package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/middlewares"
	"github.com/go-api-template/go-backend/modules/utils"
	api "github.com/go-api-template/go-backend/modules/utils/api"
	"github.com/go-api-template/go-backend/modules/utils/token"
	"github.com/go-api-template/go-backend/services"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
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
//	@Description	Create a new user. The first user created is an admin.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserSignUp	true	"User sign up"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	api.Error
//	@Failure		409		{object}	api.Error
//	@Failure		412		{object}	api.Error
//	@Failure		502		{object}	api.Error
//	@Router			/auth/signup [post]
func (c *AuthControllerImpl) SignUp(ctx *gin.Context) {
	var payload *models.UserSignUp

	// Bind the request body to the payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Check if the passwords match
	if payload.Password != payload.PasswordConfirmation {
		api.Ctx(ctx).PreconditionFailed().
			WithCode("password_match").
			WithDescription("passwords do not match").
			Send()
		return
	}

	// Sign up the user
	user, err := c.userService.Create(payload)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			api.Ctx(ctx).Conflict().WithError(err).Send()
			return
		}
		api.Ctx(ctx).BadGateway().WithError(err).Send()
		return
	}

	// Send verification code to user email address in background
	go func() {
		if user.VerificationToken != "" {
			if err := c.mailerService.SendVerificationToken(user); err != nil {
				log.Error().Err(err).Msg("Failed to send verification code")
			}
		}
	}()

	// Send the response
	api.Ctx(ctx).Created().SendRaw(user.Response())
}

// SignIn godoc
//
//	@Summary		Sign in a user
//	@Description	Sign in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			account	body		models.UserSignIn	true	"User credential"
//	@Success		201		{object}	token.AccessToken
//	@Failure		400		{object}	api.Error
//	@Failure		404		{object}	api.Error
//	@Router			/auth/signin [post]
func (c *AuthControllerImpl) SignIn(ctx *gin.Context) {
	var payload *models.UserSignIn

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Find the user by email
	user, err := c.userService.FindByEmail(payload.Email)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	if user == nil {
		api.Ctx(ctx).NotFound().
			WithCode("invalid").
			WithDescription("Invalid email or password").
			Send()
		return
	}

	// Check if the password is correct
	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		api.Ctx(ctx).BadRequest().
			WithCode("invalid").
			WithDescription("Invalid email or password").
			Send()
		return
	}
	// Check if the account is verified
	if !user.Verified {
		api.Ctx(ctx).BadRequest().
			WithCode("invalid").
			WithDescription("Account not verified").
			Send()
		return
	}

	// Sign in the user
	c.signin(ctx, user)
}

func (c *AuthControllerImpl) signin(ctx *gin.Context, user *models.User) {

	// Generate the tokens
	tokens, err := c.generateTokens(user)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Set the cookies
	ctx.SetCookie(CtxAccessToken, tokens.AccessToken, tokens.ExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, tokens.RefreshToken, tokens.RefreshExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", tokens.ExpiresIn, "/", "localhost", false, false)

	// Send the response
	api.Ctx(ctx).Created().SendRaw(&tokens)
}

// SignOut godoc
//
//	@Summary		Sign out the current user
//	@Description	Sign out the current user.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.Success
//	@Failure		401	{object}	api.Error
//	@Router			/auth/signout [get]
func (c *AuthControllerImpl) SignOut(ctx *gin.Context) {
	// Get the user from the context
	user, err := middlewares.GetUserFromContext(ctx)
	if err != nil || user == nil {
		api.Ctx(ctx).Unauthorized().
			WithCode("invalid").
			WithDescription("The user must be logged in").
			Send()
		ctx.Abort()
		return
	}

	// Sign out the user
	c.signOut(ctx, user)
}

func (c *AuthControllerImpl) signOut(ctx *gin.Context, _ *models.User) {
	// Clear the cookies
	ctx.SetCookie(CtxAccessToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, "", -1, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "", -1, "/", "localhost", false, false)

	// Send the response
	api.Ctx(ctx).Ok().Success()
}

// Welcome godoc
//
//	@Summary		Send welcome email
//	@Description	This re-sends the welcome email to the user if the user is not verified
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	body		models.UserEmail	true	"User email"	Format(email)
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	api.Error
//	@Failure		403		{object}	api.Error
//	@Failure		404		{object}	api.Error
//	@Router			/auth/welcome [post]
func (c *AuthControllerImpl) Welcome(ctx *gin.Context) {
	var payload models.UserEmail
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Get the user from the database by email
	user, err := c.userService.FindByEmail(payload.Email)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	if user == nil {
		message := "unknown user"
		api.Ctx(ctx).NotFound().
			WithCode("unknown_user").
			WithDescription(message).
			Send()
		return
	}
	if user.Verified {
		api.Ctx(ctx).Forbidden().
			WithCode("account_verified").
			WithDescription("Account already verified").
			Send()
		return
	}

	// Generate a new Verification Code
	user.VerificationToken = utils.Encode(utils.GenerateRandomString(32))

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Send verification code to user email address in background
	go func() {
		if err := c.mailerService.SendVerificationToken(user); err != nil {
			log.Error().Err(err).Msg("Failed to send the verification token")
		}
	}()

	// Send the response
	api.Ctx(ctx).Created().SendRaw(user.Response())
}

// VerifyEmail godoc
//
//	@Summary		Verify the email address
//	@Description	Verify the email address from verification code sent by email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"verification code sent by email"
//	@Success		200		{object}	api.Success
//	@Failure		400		{object}	api.Error
//	@Failure		404		{object}	api.Error
//	@Router			/auth/verify/{token} [get]
func (c *AuthControllerImpl) VerifyEmail(ctx *gin.Context) {
	// Get the verification code passed in the url
	verificationToken := ctx.Params.ByName("token")

	// Get the user with the verification code
	user, err := c.userService.FindByVerificationToken(verificationToken)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	if user == nil {
		message := "the user belonging to this code no longer exists verificationToken " + verificationToken
		api.Ctx(ctx).NotFound().
			WithCode("account_error").
			WithDescription(message).
			Send()
		return
	}

	// Set the user as verified
	user.Verified = true

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	api.Ctx(ctx).Ok().
		WithCode("email_verified").
		WithDescription("Email verified successfully").
		Send()
}

// RefreshTokens godoc
//
//	@Summary		Refresh the access token
//	@Description	Refresh the access token using the refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	body		models.UserToken	true	"Refresh token"
//	@Success		201		{object}	token.AccessToken
//	@Failure		400		{object}	api.Error
//	@Failure		403		{object}	api.Error
//	@Failure		404		{object}	api.Error
//	@Router			/auth/refresh [post]
func (c *AuthControllerImpl) RefreshTokens(ctx *gin.Context) {
	var payload *models.UserToken
	var refreshToken string

	body, _ := io.ReadAll(ctx.Request.Body)
	println(string(body))
	ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

	// Get the refresh token from the body
	if errBind := ctx.ShouldBindJSON(&payload); errBind != nil {
		// If the body is empty, get the refresh token from the cookie
		ctxToken, errCtx := ctx.Cookie(CtxRefreshToken)
		if ctxToken != "" {
			refreshToken = ctxToken
		} else if errors.Is(errCtx, http.ErrNoCookie) {
			api.Ctx(ctx).BadRequest().WithError(errBind).WithError(errCtx).Send()
			return
		}
	} else {
		refreshToken = payload.Token
	}

	// Validate the refresh token
	claims, err := token.Validate(refreshToken, config.Config.Tokens.Refresh.PublicKey)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Get the subscriber from the refresh token
	subscriber, err := claims.GetSubject()
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Get the user UUID from the refresh token
	subscriberUUID, err := uuid.Parse(fmt.Sprint(subscriber))
	if err != nil {
		api.Ctx(ctx).BadRequest().
			WithCode("token_error").
			WithDescription("cannot get user id from token").
			Send()
		return
	}

	// Get the user from the database
	user, err := c.userService.FindById(subscriberUUID)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	if user == nil {
		api.Ctx(ctx).NotFound().
			WithCode("token_error").
			WithDescription("the user belonging to this token no longer exists").
			Send()
		return
	}

	// Generate a new access token
	tokens, err := c.generateTokens(user)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Set the cookies
	ctx.SetCookie(CtxAccessToken, tokens.AccessToken, tokens.ExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxRefreshToken, tokens.RefreshToken, tokens.RefreshExpiresIn, "/", "localhost", false, true)
	ctx.SetCookie(CtxLoggedIn, "true", tokens.ExpiresIn, "/", "localhost", false, false)

	// Send the response
	api.Ctx(ctx).Created().SendRaw(&tokens)
}

// ForgotPassword godoc
//
//	@Summary		Send a reset token by email
//	@Description	Send a reset token by email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	body		models.UserEmail	true	"User email"	Format(email)
//	@Success		201		{object}	api.Success
//	@Failure		400		{object}	api.Error
//	@Failure		401		{object}	api.Error
//	@Router			/auth/forgot-password/{email} [post]
func (c *AuthControllerImpl) ForgotPassword(ctx *gin.Context) {
	var payload models.UserEmail
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Get the user from the database by email
	user, err := c.userService.FindByEmail(payload.Email)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	if user == nil {
		api.Ctx(ctx).NotFound().
			WithCode("unknown_user").
			WithDescription("unknown user").
			Send()
		return
	}
	if !user.Verified {
		api.Ctx(ctx).Unauthorized().
			WithCode("not_verified").
			WithDescription("Account not verified").
			Send()
		return
	}

	// Generate a reset token
	user.SetResetToken(utils.Encode(utils.GenerateRandomString(32)))

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Send verification code to user email address in background
	go func() {
		if err := c.mailerService.SendResetToken(user); err != nil {
			log.Error().Err(err).Msg("Failed to send the reset token")
		}
	}()

	// Send the response
	// Add the verification code if the app is in debug mode
	if config.Config.App.Debug {
		api.Ctx(ctx).Created().
			WithCode("reset_token").
			WithDescription("Reset token sent successfully").
			WithData(user.ResetToken).
			Send()
	} else {
		api.Ctx(ctx).Created().
			WithCode("reset_token").
			WithDescription("Reset token sent successfully").
			Send()
	}
}

// ResetPassword godoc
//
//	@Summary		Reset the user password
//	@Description	Reset the user password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			password	body		models.UserPasswordConfirmation	true	"New password with confirmation"
//	@Success		200			{object}	api.Success
//	@Failure		400			{object}	api.Error
//	@Failure		401			{object}	api.Error
//	@Failure		403			{object}	api.Error
//	@Failure		404			{object}	api.Error
//	@Failure		412			{object}	api.Error
//	@Router			/auth/reset-password/{token} [patch]
func (c *AuthControllerImpl) ResetPassword(ctx *gin.Context) {

	// Get the reset token from the url
	resetToken := ctx.Params.ByName("token")

	// Get the password from the body
	var payload *models.UserPasswordConfirmation
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Check if the passwords match
	if payload.Password != payload.PasswordConfirmation {
		api.Ctx(ctx).PreconditionFailed().
			WithCode("password_error").
			WithDescription("passwords do not match").
			Send()
		return
	}

	// Get the user with the reset token
	user, err := c.userService.FindByResetPasswordToken(resetToken)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	if user == nil {
		message := "the user belonging to this code no longer exists. The reset token was " + resetToken
		api.Ctx(ctx).NotFound().
			WithCode("user_error").
			WithDescription(message).
			Send()
		return
	}

	// This is the new password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	user.Password = hashedPassword

	// Clear the reset token
	user.SetResetToken("")

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
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
//	@Success		201			{object}	api.Success
//	@Failure		400			{object}	api.Error
//	@Failure		401			{object}	api.Error
//	@Failure		403			{object}	api.Error
//	@Failure		404			{object}	api.Error
//	@Failure		412			{object}	api.Error
//	@Router			/auth/change-password [post]
func (c *AuthControllerImpl) ChangePassword(ctx *gin.Context) {

	// Get the user from the context
	user, err := middlewares.GetUserFromContext(ctx)
	if err != nil || user == nil {
		api.Ctx(ctx).Unauthorized().
			WithCode("user_error").
			WithDescription("The user must be logged in").
			Send()
		ctx.Abort()
		return
	}

	// Get the new password from the body
	var payload *models.UserPasswordConfirmation
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Check if the passwords match
	if payload.Password != payload.PasswordConfirmation {
		api.Ctx(ctx).PreconditionFailed().
			WithCode("password_error").
			WithDescription("passwords do not match").
			Send()
		return
	}

	// This is the new password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}
	user.Password = hashedPassword

	// Update the user
	if _, err := c.userService.Update(user.ID, user); err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return
	}

	// Sign in the user
	c.signin(ctx, user)
}

func (c *AuthControllerImpl) generateTokens(user *models.User) (*token.AccessToken, error) {

	now := time.Now().UTC()
	issuer := config.Config.App.Name
	subject := user.ID.String()
	audience := user.Email

	// Generate access tokens
	accessToken, err := token.Create(now, config.Config.Tokens.Access.MaxAge, issuer, subject, audience, config.Config.Tokens.Access.PrivateKey)
	if err != nil {
		return nil, err
	}
	accessTokenExpiresIn := config.Config.Tokens.Access.MaxAge * 60
	accessTokenExpiresAt := now.Add(time.Duration(config.Config.Tokens.Access.MaxAge) * time.Minute)

	// Generate refresh tokens
	refreshToken, err := token.Create(now, config.Config.Tokens.Refresh.MaxAge, issuer, subject, audience, config.Config.Tokens.Refresh.PrivateKey)
	if err != nil {
		return nil, err
	}
	refreshTokenExpiresIn := config.Config.Tokens.Refresh.MaxAge * 60
	refreshTokenExpiresAt := now.Add(time.Duration(config.Config.Tokens.Refresh.MaxAge) * time.Minute)

	return &token.AccessToken{
		AccessToken:      accessToken,
		ExpiresIn:        accessTokenExpiresIn,
		ExpiresAt:        accessTokenExpiresAt.Unix(),
		TokenType:        "Bearer",
		RefreshToken:     refreshToken,
		RefreshExpiresIn: refreshTokenExpiresIn,
		RefreshExpiresAt: refreshTokenExpiresAt.Unix(),
	}, nil
}
