package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	Port          int    `json:"port"`
	EnableLogging bool   `json:"enable_logging"`
	MaxConnections int   `json:"max_connections"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func ValidateConfig(c *Config) error {
	if c.ServerAddress == "" {
		return fmt.Errorf("server_address cannot be empty")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.MaxConnections < 1 {
		return fmt.Errorf("max_connections must be at least 1")
	}
	return nil
}

func DisplayConfig(c *Config) {
	fmt.Printf("Server Configuration:\n")
	fmt.Printf("  Address: %s\n", c.ServerAddress)
	fmt.Printf("  Port: %d\n", c.Port)
	fmt.Printf("  Logging Enabled: %v\n", c.EnableLogging)
	fmt.Printf("  Max Connections: %d\n", c.MaxConnections)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <config_file>")
		os.Exit(1)
	}

	config, err := LoadConfig(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	DisplayConfig(config)
}