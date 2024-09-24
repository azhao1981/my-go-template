package config

import (
	"os"
)

type Config struct {
	Port     string
	LogLevel string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		Port:     port,
		LogLevel: logLevel,
	}, nil
}
