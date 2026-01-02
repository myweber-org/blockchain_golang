package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"min=1,max=65535"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode" default:"disable"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" default:"8080"`
	ReadTimeout  int    `yaml:"read_timeout" default:"30"`
	WriteTimeout int    `yaml:"write_timeout" default:"30"`
	DebugMode    bool   `yaml:"debug_mode"`
	LogLevel     string `yaml:"log_level" default:"info"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features []string       `yaml:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	setDefaultValues(&config)
	return &config, nil
}

func getDefaultConfigPath() string {
	paths := []string{
		"config.yaml",
		"config.yml",
		filepath.Join("config", "config.yaml"),
		filepath.Join("config", "config.yml"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	return nil
}

func setDefaultValues(config *AppConfig) {
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}
	if config.Server.LogLevel == "" {
		config.Server.LogLevel = "info"
	}
}