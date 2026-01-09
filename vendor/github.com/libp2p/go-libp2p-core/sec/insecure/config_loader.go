package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Version  string         `json:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	overrideStruct(config, "")
	return nil
}

func overrideStruct(v interface{}, prefix string) {
	// Implementation would use reflection to check struct tags
	// and override values from environment variables
	// Simplified for brevity
}

func validateConfig(config *AppConfig) error {
	var errors []string

	if config.Database.Host == "" {
		errors = append(errors, "database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		errors = append(errors, "database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		errors = append(errors, "server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		errors = append(errors, "read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		errors = append(errors, "write timeout cannot be negative")
	}
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		errors = append(errors, "invalid log level")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}
	return nil
}

func (c *AppConfig) String() string {
	masked := *c
	masked.Database.Password = "******"
	data, _ := json.MarshalIndent(masked, "", "  ")
	return string(data)
}