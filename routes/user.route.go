package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"github.com/go-api-template/go-backend/modules/middlewares"
)

type UserRoutesController struct {
	userController controllers.UserController
}

func NewUserRoutesController(userController controllers.UserController) UserRoutesController {
	return UserRoutesController{userController}
}

func (r *UserRoutesController) NewRoutes(rg *gin.RouterGroup) {
	// users routes for public users
	users := rg.Group("users")

	// users routes for authenticated and verified users
	usersVerified := users.Group("").
		Use(middlewares.VerifiedUser())
	usersVerified.GET("/me", r.userController.GetMe)
	usersVerified.PATCH("/me", r.userController.UpdateMe)
	usersVerified.DELETE("/me", r.userController.DeleteMe)

	// users routes for admin users
	usersAdmin := users.Group("").
		Use(middlewares.AdminUser())
	usersAdmin.GET("/", r.userController.List)
	usersAdmin.GET("/:id", r.userController.GetById)
	usersAdmin.PATCH("/:id", r.userController.Update)
	usersAdmin.DELETE("/:id", r.userController.Delete)
}
