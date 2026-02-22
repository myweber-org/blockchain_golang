package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL  string
    MaxConnections int
    DebugMode    bool
    AllowedHosts []string
}

func Load() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DB_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost"}, ","),
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
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
}package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	DatabaseURL  string `json:"database_url"`
	APIPort      int    `json:"api_port"`
	LogLevel     string `json:"log_level"`
	CacheEnabled bool   `json:"cache_enabled"`
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
			APIPort:      getEnvAsInt("API_PORT", 8080),
			LogLevel:     getEnv("LOG_LEVEL", "info"),
			CacheEnabled: getEnvAsBool("CACHE_ENABLED", true),
		}

		if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
			loadFromFile(configFile, instance)
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func loadFromFile(filename string, config *Config) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return
	}

	if fileConfig.DatabaseURL != "" {
		config.DatabaseURL = fileConfig.DatabaseURL
	}
	if fileConfig.APIPort != 0 {
		config.APIPort = fileConfig.APIPort
	}
	if fileConfig.LogLevel != "" {
		config.LogLevel = fileConfig.LogLevel
	}
	config.CacheEnabled = fileConfig.CacheEnabled
}package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	LogLevel   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: 8080,
		DebugMode:  false,
		LogLevel:   "info",
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.ServerPort = port
		}
	}

	if debugStr := os.Getenv("DEBUG_MODE"); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			cfg.DebugMode = debug
		}
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return cfg, nil
}