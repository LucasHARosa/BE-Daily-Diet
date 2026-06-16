package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                       string
	DatabaseURL                string
	JWTSecret                  string
	JWTAccessTokenExpiresMin   int
	JWTRefreshTokenExpiresDays int
	GeminiAPIKey               string
}

func Load() Config {
	accessMin, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRES_IN_MINUTES"))
	refreshDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRES_IN_DAYS"))
	if accessMin == 0 {
		accessMin = 30
	}
	if refreshDays == 0 {
		refreshDays = 7
	}

	return Config{
		Port:                       getEnv("PORT", "3333"),
		DatabaseURL:                os.Getenv("DATABASE_URL"),
		JWTSecret:                  os.Getenv("JWT_SECRET"),
		JWTAccessTokenExpiresMin:   accessMin,
		JWTRefreshTokenExpiresDays: refreshDays,
		GeminiAPIKey:               os.Getenv("GEMINI_API_KEY"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
