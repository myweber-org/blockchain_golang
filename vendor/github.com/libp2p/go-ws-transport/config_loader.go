package config

import (
    "fmt"
    "io/ioutil"
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

func LoadConfig(filename string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}