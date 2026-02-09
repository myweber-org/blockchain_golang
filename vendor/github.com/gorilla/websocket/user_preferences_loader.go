
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type UserPreferences struct {
	Theme      string   `json:"theme"`
	Language   string   `json:"language"`
	AutoSave   bool     `json:"auto_save"`
	FontSize   int      `json:"font_size"`
	EnabledModules []string `json:"enabled_modules"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open preferences file: %w", err)
	}
	defer file.Close()

	var prefs UserPreferences
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&prefs); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if err := validatePreferences(&prefs); err != nil {
		return nil, fmt.Errorf("preferences validation failed: %w", err)
	}

	return &prefs, nil
}

func validatePreferences(prefs *UserPreferences) error {
	if prefs.Theme == "" {
		return fmt.Errorf("theme cannot be empty")
	}
	if prefs.Language == "" {
		return fmt.Errorf("language cannot be empty")
	}
	if prefs.FontSize < 8 || prefs.FontSize > 72 {
		return fmt.Errorf("font size must be between 8 and 72")
	}
	return nil
}

func main() {
	prefs, err := LoadPreferences("config/user_preferences.json")
	if err != nil {
		log.Fatalf("Error loading preferences: %v", err)
	}

	fmt.Printf("Loaded preferences:\n")
	fmt.Printf("Theme: %s\n", prefs.Theme)
	fmt.Printf("Language: %s\n", prefs.Language)
	fmt.Printf("AutoSave: %v\n", prefs.AutoSave)
	fmt.Printf("FontSize: %d\n", prefs.FontSize)
	fmt.Printf("Enabled Modules: %v\n", prefs.EnabledModules)
}