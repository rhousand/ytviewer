package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	APIKey        string   `json:"api_key"`
	Subscriptions []string `json:"subscriptions"` // YouTube channel IDs
	MaxVideos     int64    `json:"max_videos"`
	MPVOptions    struct {
		MaxResolution  string `json:"max_resolution"`
		HardwareAccel  bool   `json:"hardware_accel"`
		CacheSize      string `json:"cache_size"`
		MarkAsWatched  bool   `json:"mark_as_watched"`
	} `json:"mpv_options"`
	CacheDuration int `json:"cache_duration"` // Cache duration in minutes
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.json")
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		return createDefaultConfig(configDir)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Set default values if not specified
	if config.MaxVideos == 0 {
		config.MaxVideos = 10
	}
	
	// Set default cache duration to 30 minutes if not specified
	if config.CacheDuration == 0 {
		config.CacheDuration = 30
	}

	return &config, nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "ytviewer")
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("error creating config directory: %w", err)
	}

	return configDir, nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configDir string) (*Config, error) {
	config := &Config{
		APIKey:        "YOUR_YOUTUBE_API_KEY",
		Subscriptions: []string{},
		MaxVideos:     10,
		MPVOptions: struct {
			MaxResolution  string `json:"max_resolution"`
			HardwareAccel  bool   `json:"hardware_accel"`
			CacheSize      string `json:"cache_size"`
			MarkAsWatched  bool   `json:"mark_as_watched"`
		}{
			MaxResolution:  "1080",
			HardwareAccel:  true,
			CacheSize:      "150M",
			MarkAsWatched:  true,
		},
		CacheDuration: 30,
	}

	// Create config file
	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error creating default config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("error writing default config: %w", err)
	}

	fmt.Printf("Created default config at %s. Please edit it to add your YouTube API key.\n", configPath)
	return config, nil
} 