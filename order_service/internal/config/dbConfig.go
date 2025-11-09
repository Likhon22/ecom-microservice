package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type DBConfig struct {
	DBUrl        string
	DBName       string
	DBDriver     string
	MaxOpenConns int
	MaxIdleConns int
}

func LoadDBConfig() *DBConfig {
	dbUrl := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	dbDriver := os.Getenv("DB_DRIVER")
	dbMaxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	dbMaxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")

	dbMaxOpenConns, err := strconv.Atoi(dbMaxOpenConnsStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid DB_MAX_OPEN_CONNS")
	}

	dbMaxIdleConns, err := strconv.Atoi(dbMaxIdleConnsStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid DB_MAX_IDLE_CONNS")
	}

	if dbUrl == "" || dbName == "" || dbDriver == "" {
		log.Fatal().Msg("DB information is required")
	}

	return &DBConfig{
		DBUrl:        dbUrl,
		MaxOpenConns: dbMaxOpenConns,
		MaxIdleConns: dbMaxIdleConns,
		DBName:       dbName,
		DBDriver:     dbDriver,
	}
}
