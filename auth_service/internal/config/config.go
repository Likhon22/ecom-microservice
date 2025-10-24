package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Version           string
	ServiceName       string
	Addr              string
	User_Service_Addr string
	DBCnf             *DBConfig
	AuthCnf           *AuthConfig
}

var (
	config *Config
	once   sync.Once
)

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("no .env file found, using system environment variables")
	}
	version := os.Getenv("VERSION")
	serviceName := os.Getenv("SERVICE_NAME")
	addr := os.Getenv("ADDR")
	user_service_addr := os.Getenv("USER_SERVICE_ADDR")

	config = &Config{
		Version:           version,
		ServiceName:       serviceName,
		Addr:              addr,
		User_Service_Addr: user_service_addr,
		DBCnf:             LoadDBConfig(),
		AuthCnf:           LoadAuthConfig(),
	}
	validateMainConfig(config)
}
func GetConfig() *Config {
	once.Do(loadConfig)
	return config

}
func validateMainConfig(cfg *Config) {
	if cfg.Version == "" || cfg.Addr == "" || cfg.ServiceName == "" || cfg.User_Service_Addr == "" {
		log.Fatal().Msg("missing core service environment variables")
	}

}
