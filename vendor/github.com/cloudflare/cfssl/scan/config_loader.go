package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKey     string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var err error

	cfg.ServerPort, err = getEnvAsInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getEnv("DB_HOST", "localhost")
	
	cfg.DBPort, err = getEnvAsInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode, err = getEnvAsBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	cfg.APIKey = getEnv("API_KEY", "")
	if cfg.APIKey == "" {
		return nil, errors.New("API_KEY environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) (int, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, errors.New("invalid integer value for " + key)
	}
	return value, nil
}

func getEnvAsBool(key string, defaultValue bool) (bool, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue, nil
	}
	valueStr = strings.ToLower(valueStr)
	if valueStr == "true" || valueStr == "1" || valueStr == "yes" {
		return true, nil
	}
	if valueStr == "false" || valueStr == "0" || valueStr == "no" {
		return false, nil
	}
	return false, errors.New("invalid boolean value for " + key)
}