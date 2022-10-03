package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateeTimeColumnName(t *testing.T) {
	tests := []struct {
		shortStepName string
		begin         bool
		want          string
	}{
		{
			shortStepName: "cnc3-qc",
			begin:         true,
			want:          "cnc3-qc.insight.test_attributes.uut_start",
		},
		{
			shortStepName: "cnc3-qc",
			begin:         false,
			want:          "cnc3-qc.insight.test_attributes.uut_stop",
		},
		{
			shortStepName: "",
			begin:         false,
			want:          "",
		},
	}

	a := require.New(t)
	for _, tt := range tests {
		a.Equal(tt.want, getTimeColumnName(tt.shortStepName, tt.begin))
	}
}
