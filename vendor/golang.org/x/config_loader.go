package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL  string `json:"database_url"`
	APIPort      int    `json:"api_port"`
	LogLevel     string `json:"log_level"`
	CacheTimeout int    `json:"cache_timeout"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
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
	return &config, nil
}

func replaceEnvVars(value string) string {
	return os.ExpandEnv(value)
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return ErrMissingDatabaseURL
	}
	if c.APIPort < 1 || c.APIPort > 65535 {
		return ErrInvalidPort
	}
	return nil
}

var (
	ErrMissingDatabaseURL = errors.New("database URL is required")
	ErrInvalidPort        = errors.New("port must be between 1 and 65535")
)