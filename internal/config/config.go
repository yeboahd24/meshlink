package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Network NetworkConfig `json:"network"`
	Media   MediaConfig   `json:"media"`
	UI      UIConfig      `json:"ui"`
}

type NetworkConfig struct {
	Port         int    `json:"port"`
	DiscoveryKey string `json:"discovery_key"`
	MaxPeers     int    `json:"max_peers"`
}

type MediaConfig struct {
	VideoCodec   string `json:"video_codec"`
	AudioCodec   string `json:"audio_codec"`
	Bitrate      int    `json:"bitrate"`
	Resolution   string `json:"resolution"`
	FrameRate    int    `json:"frame_rate"`
}

type UIConfig struct {
	Theme      string `json:"theme"`
	Fullscreen bool   `json:"fullscreen"`
	ShowStats  bool   `json:"show_stats"`
}

func DefaultConfig() *Config {
	return &Config{
		Network: NetworkConfig{
			Port:         8080,
			DiscoveryKey: "meshlink-church",
			MaxPeers:     50,
		},
		Media: MediaConfig{
			VideoCodec: "h264",
			AudioCodec: "aac",
			Bitrate:    2000,
			Resolution: "1280x720",
			FrameRate:  30,
		},
		UI: UIConfig{
			Theme:      "dark",
			Fullscreen: false,
			ShowStats:  true,
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	// Always use defaults first
	config := DefaultConfig()
	
	// Try to load custom config if it exists
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			var customConfig Config
			if err := json.Unmarshal(data, &customConfig); err == nil {
				// Merge custom config with defaults
				config = &customConfig
			}
		}
	}
	
	return config, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}