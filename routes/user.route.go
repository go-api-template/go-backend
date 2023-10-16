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
	// user routes for logged-in users
	user := rg.Group("user").
		Use(middlewares.VerifiedUser())
	user.GET("/me", r.userController.GetMe)
	user.PATCH("/me", r.userController.UpdateMe)
	user.DELETE("/me", r.userController.DeleteMe)

	// user routes for admin users
	userSecured := rg.Group("user").
		Use(middlewares.AdminUser())
	userSecured.GET("/:id", r.userController.FindById)
	userSecured.PATCH("/:id", r.userController.Update)
	userSecured.DELETE("/:id", r.userController.Delete)

	// users routes for admin users
	usersSecured := rg.Group("users").
		Use(middlewares.AdminUser())
	usersSecured.GET("/", r.userController.FindAll)
}
