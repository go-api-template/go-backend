package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"sync"
)

type Routes struct {
	// SwaggerRoutes is the controller for swagger documentation
	SwaggerRoutes SwaggerRouteController

	// HealthRoutes and InfoRoute are the controllers which handle
	// various health and status checks
	HealthRoutes HealthRouteController
	InfoRoute    StatusRouteController

	// Routes below are routes for the various controllers
	AuthRoutes AuthRoutesController
	UserRoutes UserRoutesController
}

var (
	// r is the singleton instance of the Routes struct
	r *Routes
	// Prevent multiple initialization
	once sync.Once
)

func NewRoutes(gr *gin.Engine, c *controllers.Controllers) *Routes {
	once.Do(func() {
		r = &Routes{}
		r.initialize(c)
		r.mountRoutes(gr)
	})
	return r
}

func (r *Routes) initialize(c *controllers.Controllers) {
	// Initialize routes
	r.SwaggerRoutes = NewSwaggerRouteController(c.SwaggerController)
	r.HealthRoutes = NewHealthRouteController(c.HealthController)
	r.InfoRoute = NewStatusRouteController(c.StatusController)
	r.AuthRoutes = NewAuthRoutesController(c.AuthController)
	r.UserRoutes = NewUserRoutesController(c.UserController)
}

func (r *Routes) mountRoutes(gr *gin.Engine) {
	// Routes declared here are mounted to "/" path
	rgRoot := gr.Group("/")
	r.HealthRoutes.NewRoutes(rgRoot)
	r.InfoRoute.NewRoutes(rgRoot)

	// Routes declared here are mounted to "/docs" path
	rgDocs := gr.Group("/docs")
	r.SwaggerRoutes.NewRoutes(rgDocs)

	// Routes declared here are mounted to "/api/v1" path
	rgApiV1 := gr.Group("/api/v1")
	r.AuthRoutes.NewRoutes(rgApiV1)
	r.UserRoutes.NewRoutes(rgApiV1)
}
