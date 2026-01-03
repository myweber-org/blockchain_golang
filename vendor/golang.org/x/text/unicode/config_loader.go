package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	DebugMode    bool   `json:"debug_mode"`
	LogLevel     string `json:"log_level"`
}

type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)
	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	if val := os.Getenv("DB_HOST"); val != "" {
		config.Database.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		if port, err := parseInt(val); err == nil {
			config.Database.Port = port
		}
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		if port, err := parseInt(val); err == nil {
			config.Server.Port = port
		}
	}
	if val := os.Getenv("DEBUG_MODE"); val != "" {
		config.Server.DebugMode = val == "true" || val == "1"
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.Server.LogLevel = val
	}
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}