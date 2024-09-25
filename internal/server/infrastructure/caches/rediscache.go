package caches

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(options *redis.Options) *RedisCache {
	client := redis.NewClient(options)
	return &RedisCache{client}
}
func (rc *RedisCache) AddSession(sessionId, userId string, expired time.Duration) error {
	_, err := rc.client.Set(context.Background(), "sessionID_"+sessionId, userId, expired).Result()
	return err
}
func (rc *RedisCache) DelSession(sessionId string) error {
	_, err := rc.client.Del(context.Background(), "sessionID_"+sessionId).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return ErrNoSession
		default:
			return err
		}
	}
	return nil
}
