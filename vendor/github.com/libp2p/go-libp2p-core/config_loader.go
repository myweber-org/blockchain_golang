package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty")
	}

	fullPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if c.Database.Host == "" {
		return errors.New("database host is required")
	}

	if c.Database.Name == "" {
		return errors.New("database name is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}