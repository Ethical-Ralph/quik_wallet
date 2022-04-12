package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	conn *redis.Client
}

func ConnectRedis() (*Redis, error) {
	redisClient := Redis{}

	redisConn := redis.NewClient(&redis.Options{})

	_, err := redisConn.Ping().Result()
	if err != nil {
		return nil, err
	}

	fmt.Print("Redis server connected")

	return &redisClient, nil
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.conn.Set(key, value, 0).Err()
}
