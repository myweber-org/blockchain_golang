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
}