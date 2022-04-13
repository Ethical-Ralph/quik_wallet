package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	conn *redis.Client
}

func ConnectRedis(options *redis.Options) (*Redis, error) {
	redisClient := Redis{}

	redisConn := redis.NewClient(options)

	_, err := redisConn.Ping().Result()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis server connected")

	redisClient.conn = redisConn
	return &redisClient, nil
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.conn.Set(key, value, 0).Err()
}

func (r *Redis) Get(key string) string {
	val, err := r.conn.Get(key).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		return ""
	} else {
		return val
	}
}

func (r *Redis) Del(key string) {
	r.conn.Del(key)
}
