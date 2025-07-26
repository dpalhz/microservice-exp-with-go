package redisdb

import "github.com/redis/go-redis/v9"

type Config struct {
	Addr     string
	Password string
	DB       int
}

// NewClient creates a new redis client.
func NewClient(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
