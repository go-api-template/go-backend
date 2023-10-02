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
}

var (
	// c is the singleton instance of the Controllers struct
	c *Controllers
	// Prevent multiple initialization
	once sync.Once
)

func NewControllers(s *services.Services) *Controllers {
	once.Do(func() {
		c = &Controllers{}
		c.initialize(s)
	})
	return c
}

func (c *Controllers) initialize(s *services.Services) {
	// Initialize controllers
	c.SwaggerController = NewSwaggerController()
	c.StatusController = NewStatusController(s.StatusService)
	c.HealthController = NewHealthController()
}
