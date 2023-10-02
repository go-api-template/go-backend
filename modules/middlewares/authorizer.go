package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"github.com/go-api-template/go-backend/models"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/utils"
	httputil "github.com/go-api-template/go-backend/modules/utils/http"
	"github.com/go-api-template/go-backend/services"
	"github.com/google/uuid"
	"strings"
)

type Authorizer struct {
	userService services.IUserService
}

var (
	authorizer *Authorizer
	CtxUser    = "current_user"
)

// InitializeAuthorizer initializes the authorizer
func InitializeAuthorizer(userService services.IUserService) {
	authorizer = &Authorizer{userService}
}

// contextUser extracts the user from the context
// The access token can be passed in the Authorization header or in a cookie
// The access token is validated and the user is extracted from the token
// The user is then added to the context
func (a Authorizer) contextUser(ctx *gin.Context) (*models.User, error) {
	var accessToken string

	// Get the access token from the Authorization header
	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	}

	// Get the access token from the cookie
	if accessToken == "" {
		cookie, err := ctx.Cookie(controllers.CtxAccessToken)
		if err == nil {
			accessToken = cookie
		}
	}

	// Return an error if the access token is empty
	if accessToken == "" {
		return nil, errors.New("you are not logged in")
	}

	// Get the subscriber from the access token
	subscriber, err := utils.ValidateToken(accessToken, config.Config.Tokens.Access.PublicKey)
	if err != nil {
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
			httputil.Ctx(ctx).Unauthorized().Message(err.Error())
			ctx.Abort()
			return
		}

		ctx.Set(CtxUser, user)
		ctx.Next()
	}
}
