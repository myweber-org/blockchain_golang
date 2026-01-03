package config

import (
	"errors"
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
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	// Database overrides
	if envVal := os.Getenv("DB_HOST"); envVal != "" {
		config.Database.Host = envVal
	}
	if envVal := os.Getenv("DB_PORT"); envVal != "" {
		port, err := parseInt(envVal)
		if err != nil {
			return err
		}
		config.Database.Port = port
	}
	if envVal := os.Getenv("DB_USER"); envVal != "" {
		config.Database.Username = envVal
	}
	if envVal := os.Getenv("DB_PASS"); envVal != "" {
		config.Database.Password = envVal
	}
	if envVal := os.Getenv("DB_NAME"); envVal != "" {
		config.Database.Database = envVal
	}

	// Server overrides
	if envVal := os.Getenv("SERVER_PORT"); envVal != "" {
		port, err := parseInt(envVal)
		if err != nil {
			return err
		}
		config.Server.Port = port
	}
	if envVal := os.Getenv("SERVER_READ_TIMEOUT"); envVal != "" {
		timeout, err := parseInt(envVal)
		if err != nil {
			return err
		}
		config.Server.ReadTimeout = timeout
	}
	if envVal := os.Getenv("SERVER_WRITE_TIMEOUT"); envVal != "" {
		timeout, err := parseInt(envVal)
		if err != nil {
			return err
		}
		config.Server.WriteTimeout = timeout
	}
	if envVal := os.Getenv("SERVER_DEBUG"); envVal != "" {
		config.Server.DebugMode = envVal == "true"
	}
	if envVal := os.Getenv("LOG_LEVEL"); envVal != "" {
		config.Server.LogLevel = envVal
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Database.Username == "" {
		return errors.New("database username is required")
	}
	if config.Database.Database == "" {
		return errors.New("database name is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("server read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("server write timeout cannot be negative")
	}

	return nil
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}