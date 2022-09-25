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
				"time_interval": 10
			}`,
			expectedErr: false,
			want: &ServerConfig{
				NameMap: map[string]string{
					"short_a": "full_a",
					"short_b": "full_b",
				},
				TimeInterval: 10,
			},
		},
		// Duplicate key in name_map will overwrite the previous value.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_a":"full_b"
				},
				"time_interval": 10
			}`,
			expectedErr: false,
			want: &ServerConfig{
				NameMap: map[string]string{
					"short_a": "full_b",
				},
				TimeInterval: 10,
			},
		},
		// Invalid json format.
		{
			jsonContent: `{
				"name_map":{
					"short_a":"full_a",
					"short_a":"full_b",
				},
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
				"time_interval": "ab"
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
			a.Equal(test.want, got)
		}
	}
}

func buildFilePath(base string, idx int) string {
	return filepath.Join(base, fmt.Sprintf("test-%d.json", idx))
}
