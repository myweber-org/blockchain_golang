package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	APIKey      string
	Debug       bool
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	content = os.ExpandEnv(content)

	lines := strings.Split(content, "\n")
	cfg := &Config{}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DATABASE_URL":
			cfg.DatabaseURL = value
		case "API_KEY":
			cfg.APIKey = value
		case "DEBUG":
			cfg.Debug = strings.ToLower(value) == "true"
		}
	}

	return cfg, nil
}package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

type AppConfig struct {
	Environment string         `json:"environment"`
	Debug       bool           `json:"debug"`
	Database    DatabaseConfig `json:"database"`
	Server      ServerConfig   `json:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)
	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	if env := os.Getenv("APP_ENV"); env != "" {
		config.Environment = env
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.Database.Port = atoi(dbPort)
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.Username = dbUser
	}

	if dbPass := os.Getenv("DB_PASS"); dbPass != "" {
		config.Database.Password = dbPass
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.Database = dbName
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		config.Server.Port = atoi(serverPort)
	}
}

func atoi(s string) int {
	var result int
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			break
		}
		result = result*10 + int(ch-'0')
	}
	return result
}