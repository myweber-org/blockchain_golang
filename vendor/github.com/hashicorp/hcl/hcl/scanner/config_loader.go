package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		envTag := structField.Tag.Get("env")
		defaultTag := structField.Tag.Get("default")
		
		var value string
		if envTag != "" {
			if envVal := os.Getenv(envTag); envVal != "" {
				value = envVal
			} else if defaultTag != "" {
				value = defaultTag
			}
		}
		
		if value == "" {
			continue
		}
		
		if err := setFieldValue(field, value); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}
	
	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func (c *Config) Validate() error {
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return errors.New("invalid server port")
	}
	if c.DBHost == "" {
		return errors.New("database host cannot be empty")
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return errors.New("invalid database port")
	}
	return nil
}

func (c *Config) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *Config) String() string {
	var builder strings.Builder
	builder.WriteString("Configuration:\n")
	builder.WriteString(fmt.Sprintf("  ServerPort: %d\n", c.ServerPort))
	builder.WriteString(fmt.Sprintf("  DBHost: %s\n", c.DBHost))
	builder.WriteString(fmt.Sprintf("  DBPort: %d\n", c.DBPort))
	builder.WriteString(fmt.Sprintf("  DebugMode: %v\n", c.DebugMode))
	return builder.String()
}