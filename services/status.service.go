package services

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type StatusService interface {
	IsDbConnected() bool
	GetDbConnectionStatus() string
	IsRedisConnected() bool
	GetRedisConnectionStatus() string
}

type StatusServiceImpl struct {
	ctx   context.Context
	sqlDb *sql.DB
	redis *redis.Client
}

var _ StatusService = &StatusServiceImpl{}

func NewInfoService(ctx context.Context, sqlDb *sql.DB, redis *redis.Client) StatusService {
	return &StatusServiceImpl{ctx: ctx, sqlDb: sqlDb, redis: redis}
}

func (s *StatusServiceImpl) IsDbConnected() bool {
	if err := s.sqlDb.PingContext(s.ctx); err == nil {
		return true
	}
	return false
}

func (s *StatusServiceImpl) GetDbConnectionStatus() string {
	if s.IsDbConnected() {
		return "Connected"
	}
	return "Not connected"
}

func (s *StatusServiceImpl) IsRedisConnected() bool {
	if err := s.redis.Set(s.ctx, "redis", "Connected", 0).Err(); err == nil {
		if _, err := s.redis.Get(s.ctx, "redis").Result(); err == nil {
			s.redis.Del(s.ctx, "redis")
			return true
		}
	}
	return false
}

func (s *StatusServiceImpl) GetRedisConnectionStatus() string {
	if s.IsRedisConnected() {
		return "Connected"
	}
	return "Not connected"
}
