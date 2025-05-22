package config

import (
	"encoding/json"
	"os"
)

// ServerConfig holds the HTTP server settings
type ServerConfig struct {
	Address string `json:"address"` // e.g. ":8080"
}

// HostConfig defines each KVM host's endpoint and capacity
type HostConfig struct {
	ID      string `json:"id"`      // Unique host identifier
	Address string `json:"address"` // Libvirt SSH URI or IP:port
	CPU     int    `json:"cpu"`     // Available vCPUs
	Memory  int    `json:"memory"`  // Available RAM in MB
	Storage int    `json:"storage"` // Available disk in GB
}

// TerraformConfig holds the path to the environment directory
type TerraformConfig struct {
	WorkDir string `json:"work_dir"` // terraform/environments/poc
}

// Config is the main application configuration
type Config struct {
	Server    ServerConfig    `json:"server"`
	Hosts     []HostConfig    `json:"hosts"`
	Terraform TerraformConfig `json:"terraform"`
}

// Load parses config.json from working directory
func Load() (*Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
