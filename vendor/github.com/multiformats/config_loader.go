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

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DB_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
    }
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

func getEnvAsSlice(key string, defaultValue []string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, ",")
}
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASS"`
	Database string `yaml:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	if config.Server.LogLevel == "" {
		config.Server.LogLevel = "info"
	}

	if config.Database.Host == "" {
		return errors.New("database host is required")
	}

	if config.Database.Database == "" {
		return errors.New("database name is required")
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Server.ReadTimeout < 0 {
		return errors.New("read timeout must be non-negative")
	}

	if config.Server.WriteTimeout < 0 {
		return errors.New("write timeout must be non-negative")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[config.Server.LogLevel] {
		return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
	}

	return nil
}