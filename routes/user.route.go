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

	// Public users routes
	users := rg.Group("/users")
	users.Use(middlewares.ContextUser())
	//users.Use(middlewares.VerifiedUser(userService))
	users.GET("/me", r.userController.GetMe)

	// Private users routes

	// Public user routes

	// Public user routes

}
