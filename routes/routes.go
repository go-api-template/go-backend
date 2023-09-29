package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"sync"
)

type Routes struct {
	HealthRoutes  HealthRouteController
	InfoRoute     StatusRouteController
	SwaggerRoutes SwaggerRouteController
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
	r.HealthRoutes = NewHealthRouteController(c.HealthController)
	r.InfoRoute = NewStatusRouteController(c.StatusController)
}

func (r *Routes) mountRoutes(gr *gin.Engine) {
	// Public root routes
	rgRoot := gr.Group("/")
	r.HealthRoutes.NewRoutes(rgRoot)
	r.InfoRoute.NewRoutes(rgRoot)

	// docs routes
	rgDocs := gr.Group("/docs")
	r.SwaggerRoutes.NewRoutes(rgDocs)
	_ = rgDocs
}
