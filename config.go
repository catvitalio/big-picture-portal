package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	BigPictureDisplay string `json:"bigPictureDisplay"`
	MainDisplay       string `json:"mainDisplay"`
	CheckInterval     int    `json:"checkInterval"`
}

var (
	config     Config
	configPath string
)

func initConfig() error {
	appData := os.Getenv("APPDATA")
	appDir := filepath.Join(appData, "BigPicturePortal")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}
	configPath = filepath.Join(appDir, "config.json")

	loadConfig()
	return nil
}

func loadConfig() {
	data, err := os.ReadFile(configPath)
	if err != nil {
		config = Config{
			BigPictureDisplay: "external",
			MainDisplay:       "internal",
			CheckInterval:     2000,
		}
		return
	}

	json.Unmarshal(data, &config)
}

func saveConfig() {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(configPath, data, 0644)
}
