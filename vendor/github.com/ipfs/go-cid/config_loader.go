package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	SSL  struct {
		Enabled bool   `json:"enabled"`
		Cert    string `json:"certificate"`
		Key     string `json:"key"`
	} `json:"ssl"`
	MaxConnections int `json:"max_connections"`
}

func LoadConfig(path string) (*ServerConfig, error) {
	if path == "" {
		return nil, errors.New("config path cannot be empty")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.Port == 0 {
		config.Port = 8080
	}

	if config.MaxConnections <= 0 {
		config.MaxConnections = 100
	}

	if config.SSL.Enabled && (config.SSL.Cert == "" || config.SSL.Key == "") {
		return nil, errors.New("SSL enabled but certificate or key path missing")
	}

	return &config, nil
}