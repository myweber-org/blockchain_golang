package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    AllowedHosts []string
}

func LoadConfig() (*Config, error) {
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }

    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    debugStr := getEnvWithDefault("ENABLE_DEBUG", "false")
    debug, err := strconv.ParseBool(debugStr)
    if err != nil {
        return nil, err
    }

    hostsStr := getEnvWithDefault("ALLOWED_HOSTS", "localhost,127.0.0.1")
    hosts := strings.Split(hostsStr, ",")

    return &Config{
        ServerPort:  port,
        DatabaseURL: dbURL,
        EnableDebug: debug,
        AllowedHosts: hosts,
    }, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}