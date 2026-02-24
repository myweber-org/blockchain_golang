package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	Server struct {
		Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
		Port int    `json:"port" env:"SERVER_PORT" default:"8080"`
	} `json:"server"`
	Database struct {
		Driver   string `json:"driver" env:"DB_DRIVER" default:"postgres"`
		Host     string `json:"host" env:"DB_HOST" default:"localhost"`
		Port     int    `json:"port" env:"DB_PORT" default:"5432"`
		Name     string `json:"name" env:"DB_NAME" required:"true"`
		User     string `json:"user" env:"DB_USER" required:"true"`
		Password string `json:"password" env:"DB_PASSWORD" required:"true"`
		SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	} `json:"database"`
	Logging struct {
		Level  string `json:"level" env:"LOG_LEVEL" default:"info"`
		Format string `json:"format" env:"LOG_FORMAT" default:"json"`
	} `json:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config
	
	if configPath != "" {
		if err := loadFromFile(configPath, &cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}
	
	if err := loadFromEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}
	
	if err := applyDefaults(&cfg); err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}
	
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &cfg, nil
}

func loadFromFile(path string, cfg *Config) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	
	return nil
}

func loadFromEnv(cfg *Config) error {
	return processStruct(reflect.ValueOf(cfg).Elem(), "")
}

func processStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if jsonTag := structField.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				if newPrefix != "" {
					newPrefix += "_"
				}
				newPrefix += strings.ToUpper(jsonTag)
			}
			if err := processStruct(field, newPrefix); err != nil {
				return err
			}
			continue
		}
		
		envTag := structField.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		if prefix != "" {
			envTag = prefix + "_" + envTag
		}
		
		if envValue, exists := os.LookupEnv(envTag); exists {
			if field.CanSet() {
				switch field.Kind() {
				case reflect.String:
					field.SetString(envValue)
				case reflect.Int:
					var intVal int
					if _, err := fmt.Sscanf(envValue, "%d", &intVal); err == nil {
						field.SetInt(int64(intVal))
					}
				case reflect.Bool:
					lowerVal := strings.ToLower(envValue)
					field.SetBool(lowerVal == "true" || lowerVal == "1" || lowerVal == "yes")
				}
			}
		}
	}
	
	return nil
}

func applyDefaults(cfg *Config) error {
	return applyDefaultsRecursive(reflect.ValueOf(cfg).Elem())
}

func applyDefaultsRecursive(v reflect.Value) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			if err := applyDefaultsRecursive(field); err != nil {
				return err
			}
			continue
		}
		
		if !field.IsZero() {
			continue
		}
		
		defaultTag := structField.Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		
		if field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(defaultTag)
			case reflect.Int:
				var intVal int
				if _, err := fmt.Sscanf(defaultTag, "%d", &intVal); err == nil {
					field.SetInt(int64(intVal))
				}
			case reflect.Bool:
				lowerVal := strings.ToLower(defaultTag)
				field.SetBool(lowerVal == "true" || lowerVal == "1" || lowerVal == "yes")
			}
		}
	}
	
	return nil
}

func validateConfig(cfg *Config) error {
	return validateRecursive(reflect.ValueOf(cfg).Elem())
}

func validateRecursive(v reflect.Value) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		if field.Kind() == reflect.Struct {
			if err := validateRecursive(field); err != nil {
				return err
			}
			continue
		}
		
		requiredTag := structField.Tag.Get("required")
		if requiredTag == "true" && field.IsZero() {
			jsonTag := structField.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = structField.Name
			}
			return errors.New("required field '" + jsonTag + "' is not set")
		}
	}
	
	return nil
}

func (c *Config) GetDSN() string {
	db := c.Database
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}