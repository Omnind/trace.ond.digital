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
	//  StepOrder is the list of step ordering.
	StepOrder []string `json:"step_order"`
	// StepOrderShort is the list of step ordering in short name.
	StepOrderShort []string
	// StatusColumnName is the name of the column for the status.
	StatusColumnName string `json:"status_column_name"`
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

	// We need to valid the all steps in StepOrder are in NameMap and the size of them are the same.
	if len(config.StepOrder) != len(config.NameMap) {
		return nil, errors.Errorf("the size of step serials and name map are not the same, step order: %d, name map: %d", len(config.StepOrder), len(config.NameMap))
	}

	duplicateValue := make(map[string]struct{})
	for _, v := range config.NameMap {
		if _, ok := duplicateValue[v]; ok {
			return nil, errors.Errorf("duplicate full_name in name map: %q", v)
		}
		duplicateValue[v] = struct{}{}
	}

	reverseMap := make(map[string]string)
	for k, v := range config.NameMap {
		reverseMap[v] = k
	}

	for _, step := range config.StepOrder {
		if _, ok := reverseMap[step]; !ok {
			return nil, errors.Errorf("step %q is not in name map", step)
		}
	}

	for _, v := range config.StepOrder {
		config.StepOrderShort = append(config.StepOrderShort, reverseMap[v])
	}

	return &config, nil
}
