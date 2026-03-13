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
    EnableLogging bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvOrDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("invalid DATABASE_URL format")
    }
    cfg.DatabaseURL = dbURL
    
    loggingStr := getEnvOrDefault("ENABLE_LOGGING", "true")
    cfg.EnableLogging = strings.ToLower(loggingStr) == "true"
    
    maxConnStr := getEnvOrDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil || maxConn <= 0 {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS: must be positive integer")
    }
    cfg.MaxConnections = maxConn
    
    return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}