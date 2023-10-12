package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/middlewares"
	"github.com/rs/zerolog/log"
	"sync"
)

type router struct {
	ginRouter *gin.Engine
}

var (
	// r is the main router
	r *router
	// Prevent multiple initialization
	once sync.Once
)

// NewRouter creates a new router
// It is a singleton and can be called multiple times
// but will only return the same instance
func NewRouter() *gin.Engine {
	once.Do(func() {
		r = &router{}
		r.initialize()
	})
	return r.ginRouter
}

// initialize initializes the router
func (r *router) initialize() {
	log.Debug().Msg("Initializing router...")

	// Disable gin warning if not in debug mode
	if !config.Config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new gin router
	router := gin.New()

	// Add static files
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.Use(static.Serve("/", static.LocalFile("./assets", false)))

	// Add default middleware
	router.Use(gin.Recovery())
	router.Use(middlewares.ConsoleLogger())
	router.Use(middlewares.Cors())

	// Set gin router
	r.ginRouter = router
}
