package services

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
)

// Services struct holds all services.
// All services should be declared here
// and initialized in the initialize function.
type Services struct {
	MailService MailerService
	AuthService AuthService
	UserService UserService
}

var (
	// s is the singleton instance of the Services struct
	s *Services
	// Prevent multiple initialization
	once sync.Once
)

// NewServices returns the singleton instance of the Services struct.
func NewServices(ctx context.Context, gorm *gorm.DB, sql *sql.DB, redis *redis.Client) *Services {
	once.Do(func() {
		s = &Services{}
		s.initialize(ctx, gorm, sql, redis)
	})
	return s
}

// initialize initializes all services.
func (s *Services) initialize(ctx context.Context, gorm *gorm.DB, sql *sql.DB, redis *redis.Client) {
	s.MailService, _ = NewMailerService(ctx)
	s.AuthService = NewAuthService(ctx, gorm)
	s.UserService = NewUserService(ctx, gorm)
}
