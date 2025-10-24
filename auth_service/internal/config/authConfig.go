package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type AuthConfig struct {
	Jwt_Access_Token_Secret    string
	Jwt_Refresh_Token_Secret   string
	Access_Token_Exp_Duration  time.Duration
	Refresh_Token_Exp_Duration time.Duration
}

func LoadAuthConfig() *AuthConfig {
	jwtSecret := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	jwtExpireStr := os.Getenv("ACCESS_TOKEN_EXP_DURATION")
	resetTokenExpDurationStr := os.Getenv("RESET_TOKEN_EXP_DURATION")

	jwtExpire, err := time.ParseDuration(jwtExpireStr)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid ACCESS_TOKEN_EXP_DURATION format (e.g. '1h', '30m')")
	}

	resetTokenExpDuration, err := time.ParseDuration(resetTokenExpDurationStr)
	if err != nil {
		fmt.Println(resetTokenExpDuration)
		log.Fatal().Err(err).Msg("invalid RESET_TOKEN_EXP_DURATION; must be time")
	}

	if jwtSecret == "" || jwtRefreshSecret == "" {
		log.Fatal().Msg("JWT_SECRET is required")
	}

	return &AuthConfig{
		Jwt_Access_Token_Secret:    jwtSecret,
		Access_Token_Exp_Duration:  jwtExpire,
		Jwt_Refresh_Token_Secret:   jwtRefreshSecret,
		Refresh_Token_Exp_Duration: resetTokenExpDuration,
	}
}
