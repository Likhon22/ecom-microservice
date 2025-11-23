package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type DBConfig struct {
	DBUrl        string
	RedisUrl     string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func LoadDBConfig() *DBConfig {
	dbUrl := os.Getenv("DB_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	dbMaxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	dbMaxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	dbMaxIdleTime := os.Getenv("DB_MAX_IDLE_TIME")

	dbMaxOpenConns, err := strconv.Atoi(dbMaxOpenConnsStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid DB_MAX_OPEN_CONNS")
	}

	dbMaxIdleConns, err := strconv.Atoi(dbMaxIdleConnsStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid DB_MAX_IDLE_CONNS")
	}

	if dbUrl == "" {
		log.Fatal().Msg("DB_URL is required")
	}
	if redisAddr == "" {
		log.Fatal().Msg("Redis_URL is required")
	}

	return &DBConfig{
		DBUrl:        dbUrl,
		MaxOpenConns: dbMaxOpenConns,
		MaxIdleConns: dbMaxIdleConns,
		MaxIdleTime:  dbMaxIdleTime,
		RedisUrl:     redisAddr,
	}
}
