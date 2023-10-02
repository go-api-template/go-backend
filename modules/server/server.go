package server

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-api-template/go-backend/controllers"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-api-template/go-backend/modules/database/postgres"
	redis_db "github.com/go-api-template/go-backend/modules/database/redis"
	"github.com/go-api-template/go-backend/modules/middlewares"
	"github.com/go-api-template/go-backend/modules/router"
	"github.com/go-api-template/go-backend/routes"
	"github.com/go-api-template/go-backend/services"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Server is the server struct that holds the server configuration
// It is used to initialize the server
type Server struct {
	ctx         context.Context
	gormDb      *gorm.DB
	sqlDb       *sql.DB
	redis       *redis.Client
	router      *gin.Engine
	services    *services.Services
	controllers *controllers.Controllers
	routes      *routes.Routes
}

func (s *Server) Run() (err error) {
	// Get the server configuration
	c := config.Config

	// Log the server startup
	log.Debug().Msgf("Running %s server", c.App.Name)

	// Initialize the server
	s.ctx = context.TODO()
	s.gormDb, s.sqlDb = postgres_db.NewPostgres(s.ctx)
	s.redis = redis_db.NewRedis(s.ctx)
	s.router = router.NewRouter()
	s.services = services.NewServices(s.ctx, s.gormDb, s.sqlDb, s.redis)
	s.controllers = controllers.NewControllers(s.services)
	s.routes = routes.NewRoutes(s.router, s.controllers)

	// Initialize the middlewares
	middlewares.InitializeAuthorizer(s.services.UserService)

	//
	defer func(sqlDb *sql.DB) { _ = sqlDb.Close() }(s.sqlDb)
	defer func(redis *redis.Client) { _ = redis.Close() }(s.redis)

	// Server configuration
	server := &http.Server{
		Addr:              ":8080",
		Handler:           s.router,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Info().
		Str("host", c.Server.Host).
		Str("port", c.Server.Port).
		Msgf("Starting %s Server", c.App.Name)

	if err = server.ListenAndServe(); err != nil {
		log.Fatal().
			Err(err).
			Msgf("%s Server Closed", c.App.Name)
	}

	return
}
