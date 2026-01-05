package config

import (
    "encoding/json"
    "os"
    "strings"
)

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST"`
    Port     int    `json:"port" env:"DB_PORT"`
    Username string `json:"username" env:"DB_USER"`
    Password string `json:"password" env:"DB_PASS"`
    SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type AppConfig struct {
    Debug    bool           `json:"debug" env:"APP_DEBUG"`
    LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
    Database DatabaseConfig `json:"database"`
}

func LoadConfig(path string) (*AppConfig, error) {
    file, err := os.Open(path)
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
    overrideStruct(config, "")
}

func overrideStruct(s interface{}, prefix string) {
    v := reflect.ValueOf(s).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        structField := t.Field(i)

        envKey := structField.Tag.Get("env")
        if envKey == "" {
            if field.Kind() == reflect.Struct {
                nestedPrefix := prefix
                if jsonTag := structField.Tag.Get("json"); jsonTag != "" {
                    nestedPrefix = strings.TrimSuffix(prefix+"_"+strings.ToUpper(jsonTag), "_")
                }
                overrideStruct(field.Addr().Interface(), nestedPrefix)
            }
            continue
        }

        if fullKey := strings.TrimSuffix(prefix+"_"+envKey, "_"); fullKey != "" {
            if val, exists := os.LookupEnv(fullKey); exists {
                setFieldValue(field, val)
            }
        }
    }
}

func setFieldValue(field reflect.Value, value string) {
    switch field.Kind() {
    case reflect.String:
        field.SetString(value)
    case reflect.Int, reflect.Int64:
        if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
            field.SetInt(intVal)
        }
    case reflect.Bool:
        if boolVal, err := strconv.ParseBool(value); err == nil {
            field.SetBool(boolVal)
        }
    }
}