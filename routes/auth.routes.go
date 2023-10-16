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
	auth := rg.Group("/auth")
	auth.POST("/signup", r.authController.SignUp)
	auth.POST("/signin", r.authController.SignIn)
	auth.POST("/welcome/:email", r.authController.Welcome)
	auth.GET("/verify/:token", r.authController.VerifyEmail)
	auth.GET("/refresh", r.authController.RefreshTokens)
	auth.POST("/forgot-password/:email", r.authController.ForgotPassword)
	auth.PATCH("/reset-password/:token", r.authController.ResetPassword)

	authSecured := auth.Group("").
		Use(middlewares.VerifiedUser())
	authSecured.GET("/signout", r.authController.SignOut)
	authSecured.POST("/change-password", r.authController.ChangePassword)
}
