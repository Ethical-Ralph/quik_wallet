package cache

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	conn *redis.Client
}

func ConnectRedis() *Redis {
	redisClient := Redis{}
	redisClient.conn = redis.NewClient(&redis.Options{})
	return &redisClient
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.conn.Set(key, value, 0).Err()
}
