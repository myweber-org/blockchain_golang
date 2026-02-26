package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.ServerPort, err = getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getEnvString("DB_HOST", "localhost")
	
	cfg.DBPort, err = getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode, err = getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	apiKeysStr := getEnvString("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return nil, errors.New("invalid server port range")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return nil, errors.New("invalid database port range")
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolValue, nil
	}
	return defaultValue, nil
}
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort   int
	DatabaseURL  string
	LogLevel     string
	CacheEnabled bool
	MaxWorkers   int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	portStr := getEnvOrDefault("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port

	dbURL := getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
	if strings.TrimSpace(dbURL) == "" {
		return nil, errors.New("DATABASE_URL cannot be empty")
	}
	cfg.DatabaseURL = dbURL

	cfg.LogLevel = getEnvOrDefault("LOG_LEVEL", "info")
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return nil, errors.New("invalid LOG_LEVEL value")
	}

	cacheStr := getEnvOrDefault("CACHE_ENABLED", "true")
	cacheEnabled, err := strconv.ParseBool(cacheStr)
	if err != nil {
		return nil, errors.New("invalid CACHE_ENABLED value")
	}
	cfg.CacheEnabled = cacheEnabled

	workersStr := getEnvOrDefault("MAX_WORKERS", "10")
	maxWorkers, err := strconv.Atoi(workersStr)
	if err != nil || maxWorkers < 1 || maxWorkers > 100 {
		return nil, errors.New("MAX_WORKERS must be between 1 and 100")
	}
	cfg.MaxWorkers = maxWorkers

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}