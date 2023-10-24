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
	// auth routes for public users
	auth := rg.Group("/auth")
	auth.POST("/signup", r.authController.SignUp)
	auth.POST("/signin", r.authController.SignIn)
	auth.POST("/welcome", r.authController.Welcome)
	auth.GET("/verify/:token", r.authController.VerifyEmail)
	auth.POST("/refresh", r.authController.RefreshTokens)
	auth.POST("/forgot-password", r.authController.ForgotPassword)
	auth.PATCH("/reset-password/:token", r.authController.ResetPassword)

	// auth routes for authenticated and verified users
	usersVerified := auth.Group("").
		Use(middlewares.VerifiedUser())
	usersVerified.GET("/signout", r.authController.SignOut)
	usersVerified.POST("/change-password", r.authController.ChangePassword)
}
