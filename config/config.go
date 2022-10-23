package config

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/pkg/errors"
)

type StepMeta struct {
	FullStepName  string `json:"FullStepName"`
	ShortStepName string `json:"ShortStepName"`
	StepOrder     int    `json:"StepOrder"`
	FilePath      string `json:"FilePath"`
}

// ServerConfig is the configuration for the apd.
type ServerConfig struct {
	StepMetas []StepMeta `json:"StepMeta"`
	// TimeInterval is the time interval for the apd.
	TimeInterval int `json:"TimeInterval"`
	// StatusColumnName is the name of the column for the status.
	StatusColumnName string `json:"StatusColumnName"`

	//  StepOrder is the list of step ordering.
	StepOrder []string
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

	// generate step order
	sort.Slice(config.StepMetas, func(i, j int) bool {
		return config.StepMetas[i].StepOrder < config.StepMetas[j].StepOrder
	})
	for _, step := range config.StepMetas {
		config.StepOrder = append(config.StepOrder, step.FullStepName)
	}
	return &config, nil
}
