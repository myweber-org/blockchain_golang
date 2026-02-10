package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    var err error
    
    cfg.ServerPort, err = getEnvInt("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    
    cfg.DatabaseURL = getEnvString("DATABASE_URL", "localhost:5432")
    if cfg.DatabaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }
    
    cfg.EnableDebug = getEnvBool("ENABLE_DEBUG", false)
    
    cfg.MaxConnections, err = getEnvInt("MAX_CONNECTIONS", 10)
    if err != nil {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %w", err)
    }
    
    if cfg.MaxConnections <= 0 {
        return nil, fmt.Errorf("MAX_CONNECTIONS must be positive")
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
            return 0, err
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        lowerValue := strings.ToLower(value)
        return lowerValue == "true" || lowerValue == "1" || lowerValue == "yes"
    }
    return defaultValue
}