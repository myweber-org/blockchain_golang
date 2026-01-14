package config

import (
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort    int
    DatabaseURL   string
    LogLevel      string
    CacheEnabled  bool
    MaxConnections int
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:    getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL:   getEnv("DATABASE_URL", "postgres://localhost:5432/appdb"),
        LogLevel:      getEnv("LOG_LEVEL", "info"),
        CacheEnabled:  getEnvAsBool("CACHE_ENABLED", true),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 100),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
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
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *AppConfig) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }

    if cfg.DatabaseURL == "" {
        return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
    }

    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return &ConfigError{Field: "LogLevel", Message: "invalid log level"}
    }

    if cfg.MaxConnections < 1 {
        return &ConfigError{Field: "MaxConnections", Message: "must have at least 1 connection"}
    }

    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}