package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	api "github.com/go-api-template/go-backend/modules/utils/api"
	"github.com/go-api-template/go-backend/modules/utils/token"
	"github.com/go-api-template/go-backend/services"
	"github.com/google/uuid"
	"strings"
)

type Authorizer struct {
	userService services.UserService
}

var (
	authorizer *Authorizer
	CtxUser    = "current_user"
)

// InitializeAuthorizer initializes the authorizer
func InitializeAuthorizer(userService services.UserService) {
	authorizer = &Authorizer{userService}
}

// authorization extracts the JWT from the Authorization header or from the cookie
func (a Authorizer) authorization(ctx *gin.Context) (string, error) {
	var accessToken string

	// Get the access token from the Authorization header
	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	}

	// Get the access token from the cookie
	if accessToken == "" {
		cookie, err := ctx.Cookie("access_token")
		if err == nil {
			accessToken = cookie
		}
	}

	// Return an error if the access token is empty
	if accessToken == "" {
		return "", errors.New("you are not logged in")
	}

	return accessToken, nil
}

// contextUser extracts the user from the context
// The access token can be passed in the Authorization header or in a cookie
// The access token is validated and the user is extracted from the token
// The user is then added to the context
func (a Authorizer) contextUser(ctx *gin.Context) (*models.User, error) {
	// Get the access token
	accessToken, err := a.authorization(ctx)
	if err != nil {
		return nil, err
	}

	// Validate the access token
	claims, err := token.Validate(accessToken, config.Config.Tokens.Access.PublicKey)
	if err != nil {
		return nil, err
	}

	// Get the subscriber from the refresh token
	subscriber, err := claims.GetSubject()
	if err != nil {
		api.Ctx(ctx).BadRequest().WithError(err).Send()
		return nil, err
	}

	// Convert the subscriber to a UUID
	subscriberUUID, err := uuid.Parse(fmt.Sprint(subscriber))
	if err != nil {
		return nil, errors.New("cannot get user id from token")
	}

	// Get the user from the database
	user, err := a.userService.FindById(subscriberUUID)
	if err != nil {
		return nil, errors.New("the user belonging to this token no longer exists")
	}

	return user, nil
}

// ContextUser extracts the user from the cookie or the Authorization header
// and adds it to the context
func ContextUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the user from the context
		user, err := authorizer.contextUser(ctx)
		if err != nil {
			api.Ctx(ctx).Unauthorized().WithError(err).Send()
			ctx.Abort()
			return
		}

		ctx.Set(CtxUser, user)
		ctx.Next()
	}
}

// GetUserFromContext extracts the user from the context
func GetUserFromContext(ctx *gin.Context) (*models.User, error) {
	// Get the user from the context
	user, ok := ctx.Get(CtxUser)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}

	// Cast to User model
	if _, ok := user.(*models.User); !ok {
		return nil, errors.New("user not found in context")
	}

	return user.(*models.User), nil
}

// VerifiedUser checks if the user is verified
func VerifiedUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the user from the context
		user, err := authorizer.contextUser(ctx)
		if err != nil {
			api.Ctx(ctx).Unauthorized().
				WithCode("user_error").
				WithDescription("You are not logged in").
				WithError(err).
				Send()
			ctx.Abort()
			return
		}

		// Check if the user exists
		if user == nil {
			api.Ctx(ctx).Unauthorized().
				WithCode("user_error").
				WithDescription("User not found").
				Send()
			ctx.Abort()
			return
		}

		// Check if the user is verified
		if !user.Verified {
			api.Ctx(ctx).Unauthorized().
				WithCode("user_error").
				WithDescription("Your account is not verified").
				Send()
			ctx.Abort()
			return
		}

		ctx.Set(CtxUser, user)
		ctx.Next()
	}
}

// AdminUser checks if the user is an admin
func AdminUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the user from the context
		user, err := authorizer.contextUser(ctx)
		if err != nil {
			api.Ctx(ctx).Unauthorized().
				WithError(err).
				Send()
			ctx.Abort()
			return
		}

		// Check if the user exists
		if user == nil {
			api.Ctx(ctx).Unauthorized().
				WithCode("user_error").
				WithDescription("User not found").
				Send()
			ctx.Abort()
			return
		}

		// Check if the user is an admin
		if !user.Role.IsAdmin() {
			api.Ctx(ctx).Unauthorized().
				WithCode("user_error").
				WithDescription("You are not an admin").
				Send()
			ctx.Abort()
			return
		}

		ctx.Set(CtxUser, user)
		ctx.Next()
	}
}
