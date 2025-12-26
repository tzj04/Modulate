package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string
	DBURL      string
	JWTSecret  string
}

func Load() *Config {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBURL:      getEnv("DATABASE_URL", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
