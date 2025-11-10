package infra

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once          sync.Once
	redisInstance *redis.Client
)

func ConnectRedis(addr string, db int) *redis.Client {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		redisInstance = rdb
	})
	return redisInstance

}
