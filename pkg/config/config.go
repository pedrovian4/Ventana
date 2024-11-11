package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type AppConfig struct {
	Language          string `json:"language"`
	HistoryFile       string `json:"history_file"`
	MessageDirectory  string `json:"message_directory"`
	DefaultRetryCount int    `json:"default_retry_count"`
	DefaultRetryDelay int    `json:"default_retry_delay"`
	ConfigFilePath    string `json:"config_file_path"`
	AutocompleteFile  string `json:"autocomplete_file"`
	ScriptsDirectory  string `json:"scripts_directory"`
}

func NewAppConfig() *AppConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}
	configFilePath := filepath.Join(homeDir, ".config", "ventana", "ventana.json")
	scriptsDir := filepath.Join(homeDir, ".config", "ventana", "scripts")
	messageDirectory := filepath.Join(homeDir, ".config", "ventana", "messages")

	return &AppConfig{
		Language:          "en",
		HistoryFile:       filepath.Join(filepath.Dir(configFilePath), "history.csv"),
		MessageDirectory:  messageDirectory,
		DefaultRetryCount: 3,
		DefaultRetryDelay: 5,
		ConfigFilePath:    configFilePath,
		ScriptsDirectory:  scriptsDir,
	}
}

func (c *AppConfig) Load() error {
	if _, err := os.Stat(c.ConfigFilePath); os.IsNotExist(err) {
		if err := c.Save(); err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
	} else {
		file, err := os.Open(c.ConfigFilePath)
		if err != nil {
			return fmt.Errorf("failed to open config file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(c); err != nil {
			return fmt.Errorf("failed to parse config file: %v", err)
		}
	}

	if err := c.ensureHistoryFile(); err != nil {
		return fmt.Errorf("failed to ensure history file: %v", err)
	}

	if err := c.InitializeScripts(); err != nil {
		return fmt.Errorf("failed to init scripts: %v", err)
	}

	return nil
}

func (c *AppConfig) Save() error {
	configDir := filepath.Dir(c.ConfigFilePath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	file, err := os.Create(c.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("error creating config file at %s: %v", c.ConfigFilePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("error encoding config to JSON: %v", err)
	}

	return nil
}

func (c *AppConfig) Update(fields map[string]interface{}) error {
	for key, value := range fields {
		switch key {
		case "Language":
			c.Language = value.(string)
		case "HistoryFile":
			c.HistoryFile = value.(string)
		case "MessageDirectory":
			c.MessageDirectory = value.(string)
		case "DefaultRetryCount":
			c.DefaultRetryCount = value.(int)
		case "DefaultRetryDelay":
			c.DefaultRetryDelay = value.(int)
		}
	}
	return c.Save()
}

func (c *AppConfig) ensureHistoryFile() error {
	dir := filepath.Dir(c.HistoryFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating history directory: %v", err)
	}
	if _, err := os.Stat(c.HistoryFile); os.IsNotExist(err) {
		file, err := os.Create(c.HistoryFile)
		if err != nil {
			return fmt.Errorf("error creating history file: %v", err)
		}
		defer file.Close()
	}
	return nil
}

func (c *AppConfig) InitializeAutocomplete() error {
	if _, err := os.Stat(c.AutocompleteFile); os.IsNotExist(err) {
		return nil
	}
	cmd := exec.Command("bash", "-c", "source "+c.AutocompleteFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load autocomplete script: %v", err)
	}

	return nil
}

func (c *AppConfig) InitializeScripts() error {
	fmt.Println(c.ScriptsDirectory)
	if _, err := os.Stat(c.ScriptsDirectory); os.IsNotExist(err) {
		return os.MkdirAll(c.ScriptsDirectory, 0755)
	}
	return nil
}
