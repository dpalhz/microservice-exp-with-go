package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisSessionStore stores refresh tokens in Redis.
type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{client: client}
}

// Store saves the refresh token with a TTL of 7 days.
func (s *RedisSessionStore) Store(ctx context.Context, token string, userID uuid.UUID) error {
	return s.client.Set(ctx, token, userID.String(), 7*24*time.Hour).Err()
}

// GetUserID retrieves the user ID for a refresh token.
func (s *RedisSessionStore) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	val, err := s.client.Get(ctx, token).Result()
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(val)
}
