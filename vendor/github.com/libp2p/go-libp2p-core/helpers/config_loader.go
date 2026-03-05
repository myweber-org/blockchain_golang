package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DatabaseURL string `json:"database_url"`
	LogLevel string `json:"log_level"`
	CacheEnabled bool `json:"cache_enabled"`
}

var (
	instance *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		configFile := os.Getenv("CONFIG_FILE")
		if configFile == "" {
			configFile = "config.json"
		}

		file, err := os.Open(configFile)
		if err != nil {
			instance = &Config{
				ServerPort:   getEnvOrDefault("SERVER_PORT", "8080"),
				DatabaseURL:  getEnvOrDefault("DATABASE_URL", "postgres://localhost:5432/app"),
				LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
				CacheEnabled: getEnvBoolOrDefault("CACHE_ENABLED", true),
			}
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&instance); err != nil {
			panic("Failed to decode config file: " + err.Error())
		}

		overrideFromEnv(instance)
	})
	return instance
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func overrideFromEnv(cfg *Config) {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.ServerPort = port
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if cache := os.Getenv("CACHE_ENABLED"); cache != "" {
		cfg.CacheEnabled = cache == "true" || cache == "1"
	}
}