package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	DebugMode    bool   `yaml:"debug_mode"`
	LogLevel     string `yaml:"log_level"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("invalid server port")
	}

	if config.Database.Host == "" {
		return errors.New("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("invalid database port")
	}

	if config.Database.Database == "" {
		return errors.New("database name is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}