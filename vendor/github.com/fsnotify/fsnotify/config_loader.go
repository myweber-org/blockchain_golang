package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level"`
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

	overrideFromEnv(config)

	return config, nil
}

func overrideFromEnv(config *Config) {
	if val := os.Getenv("SERVER_HOST"); val != "" {
		config.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		config.Server.Port = parseInt(val)
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		config.Database.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		config.Database.Port = parseInt(val)
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		config.Database.Name = val
	}
	if val := os.Getenv("DB_USERNAME"); val != "" {
		config.Database.Username = val
	}
	if val := os.Getenv("DB_PASSWORD"); val != "" {
		config.Database.Password = val
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.LogLevel = strings.ToUpper(val)
	}
}

func parseInt(s string) int {
	var result int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		} else {
			break
		}
	}
	return result
}