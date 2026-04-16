package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:     os.Getenv("SERVER_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("FATAL: DATABASE_URL not found")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("FATAL: JWT_SECRET not found")
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8081"
	}

	return cfg
}
