package db

import "github.com/redis/go-redis/v9"

var redisInstance *redis.Client

func ConnectRedis(addr string, db int) *redis.Client {
	if redisInstance != nil {
		return redisInstance

	}
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	redisInstance = rdb
	return rdb

}
