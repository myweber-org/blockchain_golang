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
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Port        int
    DatabaseURL string
    LogLevel    string
    CacheSize   int
}

func Load() (*Config, error) {
    cfg := &Config{
        Port:        getEnvAsInt("APP_PORT", 8080),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        CacheSize:   getEnvAsInt("CACHE_SIZE", 100),
    }

    if err := validate(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }

    value, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return value
}

func validate(cfg *Config) error {
    if cfg.Port < 1 || cfg.Port > 65535 {
        return fmt.Errorf("invalid port number: %d", cfg.Port)
    }

    if !strings.HasPrefix(cfg.DatabaseURL, "postgres://") {
        return fmt.Errorf("invalid database URL format")
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
    }

    if cfg.CacheSize < 0 {
        return fmt.Errorf("cache size cannot be negative")
    }

    return nil
}