package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
		Name     string `yaml:"name" env:"DB_NAME"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.overrideFromEnv()

	return &config, nil
}

func (c *Config) overrideFromEnv() {
	c.Server.Host = getEnvOrDefault("SERVER_HOST", c.Server.Host)
	c.Server.Port = getEnvIntOrDefault("SERVER_PORT", c.Server.Port)
	
	c.Database.Host = getEnvOrDefault("DB_HOST", c.Database.Host)
	c.Database.Port = getEnvIntOrDefault("DB_PORT", c.Database.Port)
	c.Database.Username = getEnvOrDefault("DB_USER", c.Database.Username)
	c.Database.Password = getEnvOrDefault("DB_PASS", c.Database.Password)
	c.Database.Name = getEnvOrDefault("DB_NAME", c.Database.Name)
	
	c.Logging.Level = getEnvOrDefault("LOG_LEVEL", c.Logging.Level)
	c.Logging.Output = getEnvOrDefault("LOG_OUTPUT", c.Logging.Output)
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