package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    AllowedIPs []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DBHost = getEnv("DB_HOST", "localhost")
    
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, err
    }
    cfg.DBPort = dbPort
    
    debug, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
    if err != nil {
        return nil, err
    }
    cfg.DebugMode = debug
    
    ips := getEnv("ALLOWED_IPS", "127.0.0.1,::1")
    cfg.AllowedIPs = strings.Split(ips, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return strconv.ErrRange
    }
    
    if cfg.DBPort < 1 || cfg.DBPort > 65535 {
        return strconv.ErrRange
    }
    
    if cfg.DBHost == "" {
        return os.ErrInvalid
    }
    
    return nil
}