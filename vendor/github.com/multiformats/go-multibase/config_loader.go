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
    CacheEnabled bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("invalid DATABASE_URL format")
    }
    cfg.DatabaseURL = dbURL
    
    cacheEnabled := getEnvWithDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabled) == "true"
    
    maxConnStr := getEnvWithDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil || maxConn <= 0 {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS value")
    }
    cfg.MaxConnections = maxConn
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}