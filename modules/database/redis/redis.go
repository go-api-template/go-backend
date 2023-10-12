package redis_db

import (
	"context"
	"fmt"
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"sync"
)

type redisDb struct {
	redis *redis.Client
}

var (
	// r is the singleton instance of redisDb client
	r *redisDb
	// Prevent multiple initialization
	once sync.Once
)

func NewRedis(ctx context.Context) *redis.Client {
	once.Do(func() {
		r = &redisDb{}
		r.initialize(ctx)
	})
	return r.redis
}

func (r *redisDb) initialize(ctx context.Context) {
	log.Debug().Msg("Initializing Redis...")

	// Initialize redisDb
	rc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
	})

	// Validate connection to redisDb by pinging it
	err := rc.Ping(ctx).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Pinging Redis instance")
	}

	// Set redisDb
	r.redis = rc
}
