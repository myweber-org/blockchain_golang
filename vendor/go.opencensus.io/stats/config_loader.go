package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASS"`
	SSLMode  string `env:"DB_SSL_MODE" default:"disable"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := loadStruct(cfg, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

func loadStruct(v interface{}, prefix string) error {
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Struct {
			err := loadStruct(field.Addr().Interface(), prefix+fieldType.Name+"_")
			if err != nil {
				return err
			}
			continue
		}

		envKey := fieldType.Tag.Get("env")
		if envKey == "" {
			continue
		}

		if prefix != "" {
			envKey = prefix + envKey
		}

		defaultValue := fieldType.Tag.Get("default")
		envValue := os.Getenv(envKey)

		if envValue == "" {
			envValue = defaultValue
		}

		if envValue == "" && isRequired(fieldType) {
			return fmt.Errorf("required environment variable %s is not set", envKey)
		}

		if envValue != "" {
			err := setFieldValue(field, envValue)
			if err != nil {
				return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
			}
		}
	}
	return nil
}

func isRequired(field reflect.StructField) bool {
	requiredTag := field.Tag.Get("required")
	return requiredTag == "true"
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
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			slice := strings.Split(value, ",")
			field.Set(reflect.ValueOf(slice))
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}