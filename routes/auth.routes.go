package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"github.com/go-api-template/go-backend/modules/middlewares"
)

type AuthRoutesController struct {
	authController controllers.AuthController
}

func NewAuthRoutesController(authController controllers.AuthController) AuthRoutesController {
	return AuthRoutesController{authController}
}

func (r *AuthRoutesController) NewRoutes(rg *gin.RouterGroup) {

	// Public auth routes
	publicRoutes := rg.Group("/auth")
	publicRoutes.POST("/register", r.authController.UserSignUp)
	publicRoutes.GET("/welcome/:email", r.authController.Welcome)
	publicRoutes.GET("/verify/:verification_code", r.authController.VerifyEmail)
	publicRoutes.POST("/login", r.authController.UserSignIn)
	publicRoutes.GET("/refresh", r.authController.RefreshAccessToken)
	publicRoutes.GET("/password/forgot/:email", r.authController.ForgotPassword)
	publicRoutes.PATCH("/password/reset/:reset_token", r.authController.ResetPassword)

	// Private auth routes
	privateRoutes := rg.Group("/auth")
	privateRoutes.Use(middlewares.ContextUser())
	privateRoutes.Use(middlewares.VerifiedUser())
	privateRoutes.GET("/logout", r.authController.UserSignOut)
	privateRoutes.PATCH("/password/change", r.authController.ChangePassword)
}
