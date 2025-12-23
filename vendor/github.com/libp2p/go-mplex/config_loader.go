package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	portStr := os.Getenv("SERVER_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		cfg.ServerPort = port
	} else {
		cfg.ServerPort = 8080
	}

	cfg.DBHost = os.Getenv("DB_HOST")
	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr != "" {
		dbPort, err := strconv.Atoi(dbPortStr)
		if err != nil {
			return nil, err
		}
		cfg.DBPort = dbPort
	} else {
		cfg.DBPort = 5432
	}

	debugStr := os.Getenv("DEBUG_MODE")
	if debugStr != "" {
		debug, err := strconv.ParseBool(debugStr)
		if err != nil {
			return nil, err
		}
		cfg.DebugMode = debug
	} else {
		cfg.DebugMode = false
	}

	return cfg, nil
}