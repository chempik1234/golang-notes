package redis

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisStorage struct {
	db *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{db: client}
}

func (r RedisStorage) Get(key string) ([]byte, error) {
	return r.db.Get(key).Bytes()
}

func (r RedisStorage) Set(key string, val []byte, exp time.Duration) error {
	return r.db.Set(key, string(val), exp).Err()
}

func (r RedisStorage) Delete(key string) error {
	return r.db.Del(key).Err()
}

func (r RedisStorage) Reset() error {
	return r.db.FlushDB().Err()
}

func (r RedisStorage) Close() error {
	return r.db.Close()
}
