package config

import (
	"fmt"
	"os"
)

type Config struct {
	BindAddr    string
	LogLevel    string
	LogFormat   string
	DatabaseURL string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:  ":8080",
		LogLevel:  "debug",
		LogFormat: "json",
	}
}

func (c *Config) LoadFromEnv() {
	if port := os.Getenv("APP_PORT"); port != "" {
		c.BindAddr = ":" + port
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.LogLevel = logLevel
	}

	if logFormat := os.Getenv("LOG_FORMAT"); logFormat != "" {
		c.LogFormat = logFormat
	}

	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		user := get("POSTGRES_USER", "postgres")
		password := get("POSTGRES_PASSWORD", "postgres")
		dbname := get("POSTGRES_DB", "avito_tech")
		port := get("POSTGRES_PORT", "5432")
		c.DatabaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}
}

func get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
