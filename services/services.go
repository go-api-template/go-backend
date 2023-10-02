package services

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
)

type Services struct {
	StatusService StatusService
}

var (
	// s is the singleton instance of the Services struct
	s *Services
	// Prevent multiple initialization
	once sync.Once
)

func NewServices(ctx context.Context, gorm *gorm.DB, sql *sql.DB, redis *redis.Client) *Services {
	once.Do(func() {
		s = &Services{}
		s.initialize(ctx, gorm, sql, redis)
	})
	return s
}

func (s *Services) initialize(ctx context.Context, _ *gorm.DB, sql *sql.DB, redis *redis.Client) {
	s.StatusService = NewInfoService(ctx, sql, redis)
}
