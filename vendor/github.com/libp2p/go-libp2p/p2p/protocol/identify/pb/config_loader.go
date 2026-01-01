package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string `json:"database_url"`
	APIPort     int    `json:"api_port"`
	DebugMode   bool   `json:"debug_mode"`
	LogLevel    string `json:"log_level"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.DatabaseURL = replaceEnvVars(config.DatabaseURL)
	config.LogLevel = replaceEnvVars(config.LogLevel)

	return &config, nil
}

func replaceEnvVars(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envVar := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
		if envValue, exists := os.LookupEnv(envVar); exists {
			return envValue
		}
	}
	return value
}