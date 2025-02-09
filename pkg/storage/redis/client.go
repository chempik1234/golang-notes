package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2/log"
)

func NewRedisClient(redisUrl string) (*redis.Client, error) {
	option, err := redis.ParseURL(redisUrl)
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
