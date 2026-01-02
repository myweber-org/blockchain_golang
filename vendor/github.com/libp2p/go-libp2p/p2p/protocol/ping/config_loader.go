package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
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
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func DefaultConfig() *AppConfig {
	config := &AppConfig{}
	config.Server.Host = "localhost"
	config.Server.Port = 8080
	config.Database.Host = "localhost"
	config.Database.Name = "appdb"
	config.Logging.Level = "info"
	config.Logging.Output = "stdout"
	return config
}

func ValidateConfig(config *AppConfig) bool {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		log.Printf("Invalid port number: %d", config.Server.Port)
		return false
	}
	if config.Database.Name == "" {
		log.Print("Database name cannot be empty")
		return false
	}
	return true
}