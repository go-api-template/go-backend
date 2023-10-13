package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"github.com/go-api-template/go-backend/modules/config"
	common_routes "github.com/go-api-template/go-backend/routes/common"
	"sync"
)

type Routes struct {
	// SwaggerRoutes is the controller for swagger documentation
	SwaggerRoutes common_routes.SwaggerRouteController

	// HealthCheckRoutes and StatusRoutes are the controllers which handle
	// various health and status checks
	PingRoutes        common_routes.PingRoutesController
	HealthCheckRoutes common_routes.HealthCheckRouteController
	StatusRoutes      common_routes.StatusRouteController

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
	r.SwaggerRoutes = common_routes.NewSwaggerRouteController()
	r.PingRoutes = common_routes.NewPingRoutesController()
	r.HealthCheckRoutes = common_routes.NewHealthCheckRoutesController()
	r.StatusRoutes = common_routes.NewStatusRoutesController()
	r.AuthRoutes = NewAuthRoutesController(c.AuthController)
	r.UserRoutes = NewUserRoutesController(c.UserController)
}

func (r *Routes) mountRoutes(gr *gin.Engine) {
	// Define the root path
	// By default, the root path is "/"
	base := gr.Group("/")
	// However, if a base path is defined in the config file,
	// the root path will be changed to "/{base_path}"
	if config.Config.Server.BasePath != "" {
		base = gr.Group(fmt.Sprint("/", config.Config.Server.BasePath))
	}

	// Routes declared here are mounted to "/docs" path
	docs := base.Group("/docs")
	r.SwaggerRoutes.NewRoutes(docs)

	// Add routes to the base path
	r.PingRoutes.NewRoutes(base)
	r.HealthCheckRoutes.NewRoutes(base)
	r.StatusRoutes.NewRoutes(base)
	r.AuthRoutes.NewRoutes(base)
	r.UserRoutes.NewRoutes(base)
}
