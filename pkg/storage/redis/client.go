package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2/log"
)

// NewRedisClient try to connect to Redis and get the client
func NewRedisClient(redisURL string) (*redis.Client, error) {
	option, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(option)
	err = client.Ping().Err()
	if err != nil {
		return nil, err
	}
	log.Info("connected to the redis database")
	return client, nil
}
