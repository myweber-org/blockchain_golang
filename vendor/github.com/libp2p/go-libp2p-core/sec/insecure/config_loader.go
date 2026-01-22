package config

import (
	"io/ioutil"
	"log"

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

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfigWithDefaults(filename string) *Config {
	config, err := LoadConfig(filename)
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		return DefaultConfig()
	}
	return config
}

func DefaultConfig() *Config {
	var config Config
	config.Server.Host = "localhost"
	config.Server.Port = 8080
	config.Database.Host = "localhost"
	config.Database.Name = "appdb"
	config.LogLevel = "info"
	return &config
}package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Server struct {
		Port    int    `yaml:"port"`
		Host    string `yaml:"host"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		log.Printf("Invalid server port: %d", config.Server.Port)
		return false
	}

	if config.Database.Host == "" {
		log.Print("Database host cannot be empty")
		return false
	}

	if config.Logging.Level != "debug" && config.Logging.Level != "info" &&
		config.Logging.Level != "warn" && config.Logging.Level != "error" {
		log.Printf("Invalid logging level: %s", config.Logging.Level)
		return false
	}

	return true
}