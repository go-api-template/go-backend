package controllers

import (
	"github.com/go-api-template/go-backend/services"
	"sync"
)

// Controllers struct holds all controllers.
// All controllers should be declared here
// and initialized in the initialize function.
type Controllers struct {
	SwaggerController SwaggerController
	StatusController  StatusController
	HealthController  HealthController
	AuthController    AuthController
	UserController    UserController
}

var (
	// c is the singleton instance of the Controllers struct
	c *Controllers
	// Prevent multiple initialization
	once sync.Once
)

// NewControllers returns the singleton instance of the Controllers struct.
func NewControllers(s *services.Services) *Controllers {
	once.Do(func() {
		c = &Controllers{}
		c.initialize(s)
	})
	return c
}

// initialize initializes all controllers.
func (c *Controllers) initialize(s *services.Services) {
	c.SwaggerController = NewSwaggerController()
	c.StatusController = NewStatusController(s.StatusService)
	c.HealthController = NewHealthController()
	c.AuthController = NewAuthController(s.UserService, s.MailService)
	c.UserController = NewUserController(s.UserService)
}
