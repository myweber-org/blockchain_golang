package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort   int
	DatabaseURL  string
	CacheEnabled bool
	LogLevel     string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	cfg.DatabaseURL = dbURL

	cacheStr := os.Getenv("CACHE_ENABLED")
	cacheEnabled := strings.ToLower(cacheStr) == "true"
	cfg.CacheEnabled = cacheEnabled

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[strings.ToLower(logLevel)] {
		return nil, errors.New("invalid LOG_LEVEL value")
	}
	cfg.LogLevel = strings.ToLower(logLevel)

	return cfg, nil
}