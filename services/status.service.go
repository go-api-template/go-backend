package services

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
)

// IStatusService is the interface for status service
// It declares the methods that the StatusService must implement
type IStatusService interface {
	IsDbConnected() bool
	GetDbConnectionStatus() string
	IsRedisConnected() bool
	GetRedisConnectionStatus() string
}

// StatusService is the service for status
// It implements the IStatusService interface
type StatusService struct {
	ctx   context.Context
	sqlDb *sql.DB
	redis *redis.Client
}

// StatusService implements the IStatusService interface
var _ IStatusService = &StatusService{}

func NewInfoService(ctx context.Context, sqlDb *sql.DB, redis *redis.Client) IStatusService {
	return &StatusService{ctx: ctx, sqlDb: sqlDb, redis: redis}
}

func (s *StatusService) IsDbConnected() bool {
	if err := s.sqlDb.PingContext(s.ctx); err == nil {
		return true
	}
	return false
}

func (s *StatusService) GetDbConnectionStatus() string {
	if s.IsDbConnected() {
		return "Connected"
	}
	return "Not connected"
}

func (s *StatusService) IsRedisConnected() bool {
	if err := s.redis.Set(s.ctx, "redis", "Connected", 0).Err(); err == nil {
		if _, err := s.redis.Get(s.ctx, "redis").Result(); err == nil {
			s.redis.Del(s.ctx, "redis")
			return true
		}
	}
	return false
}

func (s *StatusService) GetRedisConnectionStatus() string {
	if s.IsRedisConnected() {
		return "Connected"
	}
	return "Not connected"
}
