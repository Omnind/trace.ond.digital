package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromFile(t *testing.T) {
	tests := []struct {
		jsonContent string
		expectedErr bool
		want        *ServerConfig
	}{
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_b":"full_b"
				},
				"step_order": [
					"full_a",
					"full_b"
				],
				"time_interval": 10
			}`,
			expectedErr: false,
			want: &ServerConfig{
				NameMap: map[string]string{
					"short_a": "full_a",
					"short_b": "full_b",
				},
				StepOrder:      []string{"full_a", "full_b"},
				StepOrderShort: []string{"short_a", "short_b"},
				TimeInterval:   10,
			},
		},
		// Duplicate key in name_map will overwrite the previous value.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_a":"full_b"
				},
				"step_order": [
					"full_b"
				],
				"time_interval": 10
			}`,
			expectedErr: false,
			want: &ServerConfig{
				NameMap: map[string]string{
					"short_a": "full_b",
				},
				StepOrder:      []string{"full_b"},
				StepOrderShort: []string{"short_a"},
				TimeInterval:   10,
			},
		},
		// Invalid json format.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_a":"full_b",
				},
				"step_order": [
					"full_a, full_b"
				],
				"time_interval": 10
			}`,
			expectedErr: true,
			want:        nil,
		},
		// TimeInterval is not a number but can convert to a number.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_a":"full_b"
				},
				"step_order": [
					"full_a"
				],
				"time_interval": "30"
			}`,
			expectedErr: true,
			want:        nil,
		},
		// TimeInterval is not a number and cannot convert to a number.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_a":"full_b"
				},
				"step_order": [
					"full_a"
				],
				"time_interval": "ab"
			}`,
			expectedErr: true,
			want:        nil,
		},
		// Duplicate value in namp_map.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_b":"full_a"
				},
				"step_order": [
					"full_a"
				],
				"time_interval": 10
			}`,
			expectedErr: true,
			want:        nil,
		},
		// Step in step_order is not in name_map.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_b":"full_b"
				},
				"step_order": [
					"full_a", "full_b", "full_c"
				],
				"time_interval": 10
			}`,
			expectedErr: true,
			want:        nil,
		},
		// Step in name_map is not in step_order.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_b":"full_b"
				},
				"step_order": [
					"full_a"
				],
				"time_interval": 10
			}`,
			expectedErr: true,
			want:        nil,
		},
	}
	a := require.New(t)
	tmpDir := t.TempDir()
	for idx, test := range tests {
		fp := buildFilePath(tmpDir, idx)
		err := os.WriteFile(fp, []byte(test.jsonContent), 0644)
		a.NoError(err)

		got, err := NewConfigFromFile(fp)
		if test.expectedErr {
			a.Error(err)
		} else {
			a.NoErrorf(err, "test: %d", idx)
			a.Equalf(test.want, got, "test: %d", idx)
		}
	}
}

func buildFilePath(base string, idx int) string {
	return filepath.Join(base, fmt.Sprintf("test-%d.json", idx))
}
