package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// ServerConfig is the configuration for the apd.
type ServerConfig struct {
	// NameMap is the map of mapping short name to full name.
	NameMap map[string]string `json:"name_map"`
	// TimeInterval is the time interval for the apd.
	TimeInterval int `json:"time_interval"`
}

// NewConfigFromFile creates a new ServerConfig from a file.
func NewConfigFromFile(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %q", path)
	}

	var config ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarschal config file %q", path)
	}

	return &config, nil
}
