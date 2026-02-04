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
}package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.ServerPort = getEnvAsInt("SERVER_PORT", 8080)
	cfg.DebugMode = getEnvAsBool("DEBUG_MODE", false)
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
	cfg.AllowedHosts = getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}, ",")

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

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}