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
	users := rg.Group("users").
		Use(middlewares.ContextUser()).
		Use(middlewares.VerifiedUser())
	users.GET("/me", r.userController.GetMe)
	users.PATCH("/me", r.userController.UpdateMe)
	users.DELETE("/me", r.userController.DeleteMe)
}
