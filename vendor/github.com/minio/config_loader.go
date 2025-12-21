package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) loadFromEnv() error {
	envVars := map[string]*string{
		"SERVER_HOST":     &c.Server.Host,
		"DB_HOST":         &c.Database.Host,
		"DB_NAME":         &c.Database.Name,
		"DB_USER":         &c.Database.User,
		"DB_PASSWORD":     &c.Database.Password,
		"LOG_LEVEL":       &c.LogLevel,
	}

	for envName, fieldPtr := range envVars {
		if val, exists := os.LookupEnv(envName); exists && val != "" {
			*fieldPtr = val
		}
	}

	if c.Server.Port == 0 {
		if portStr, exists := os.LookupEnv("SERVER_PORT"); exists && portStr != "" {
			port, err := parseInt(portStr)
			if err != nil {
				return err
			}
			c.Server.Port = port
		}
	}

	if c.Database.Port == 0 {
		if portStr, exists := os.LookupEnv("DB_PORT"); exists && portStr != "" {
			port, err := parseInt(portStr)
			if err != nil {
				return err
			}
			c.Database.Port = port
		}
	}

	return nil
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, errors.New("invalid integer value")
	}
	return result, nil
}