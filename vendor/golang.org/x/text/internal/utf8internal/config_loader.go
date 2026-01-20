package config

import (
	"encoding/json"
	"os"
	"sync"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

var (
	config     *AppConfig
	configOnce sync.Once
)

func LoadConfig() *AppConfig {
	configOnce.Do(func() {
		configFile := os.Getenv("CONFIG_FILE")
		if configFile == "" {
			configFile = "config.json"
		}

		file, err := os.Open(configFile)
		if err != nil {
			config = loadFromEnv()
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			config = loadFromEnv()
		}
	})
	return config
}

func loadFromEnv() *AppConfig {
	return &AppConfig{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	var result int
	if _, err := fmt.Sscanf(valueStr, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}