package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port    string `yaml:"port" env:"SERVER_PORT"`
		Timeout int    `yaml:"timeout" env:"SERVER_TIMEOUT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(config)

	return config, nil
}

func overrideWithEnvVars(config *Config) {
	config.Server.Port = getEnvOrDefault("SERVER_PORT", config.Server.Port)
	config.Server.Timeout = getEnvIntOrDefault("SERVER_TIMEOUT", config.Server.Timeout)
	config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
	config.Database.Port = getEnvOrDefault("DB_PORT", config.Database.Port)
	config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)
	config.Database.Username = getEnvOrDefault("DB_USER", config.Database.Username)
	config.Database.Password = getEnvOrDefault("DB_PASS", config.Database.Password)
	config.LogLevel = getEnvOrDefault("LOG_LEVEL", config.LogLevel)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return errors.New("server port is required")
	}
	if c.Database.Host == "" || c.Database.Name == "" {
		return errors.New("database host and name are required")
	}
	if !isValidLogLevel(c.LogLevel) {
		return errors.New("invalid log level")
	}
	return nil
}

func isValidLogLevel(level string) bool {
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	level = strings.ToLower(level)
	for _, valid := range validLevels {
		if level == valid {
			return true
		}
	}
	return false
}