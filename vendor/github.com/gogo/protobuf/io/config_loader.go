package config

import (
	"encoding/json"
	"os"
	"sync"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

var (
	config     *AppConfig
	configOnce sync.Once
)

func LoadConfig() *AppConfig {
	configOnce.Do(func() {
		config = &AppConfig{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DebugMode:  getEnv("DEBUG_MODE", "false") == "true",
		}

		if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
			loadFromFile(configFile, config)
		}
	})
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadFromFile(filename string, cfg *AppConfig) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(cfg)
}