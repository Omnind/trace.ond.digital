package worker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaster(t *testing.T) {
	tests := []struct {
		stepOrdering         []string
		nItems               int
		nInterval            int
		nHouesInOneInterval  int
		generateSerialNumber func(i int) string
		// generateStepTime is a function to generate the step time. It's responsible to generate the final step end time.
		setStepsTime func(stepsOrdering []string, item *PartItem)
		want         *ResultSet
	}{
		{
			stepOrdering:         []string{"STEP-A", "STEP-B", "STEP-C"},
			nItems:               5,
			nInterval:            7,
			nHouesInOneInterval:  24,
			generateSerialNumber: commonGenerateSerialNumber,
			setStepsTime:         commonSetStepsTime,
			want: &ResultSet{
				StepsTimeNumber: map[string]ToStepIntervalSpent{
					"STEP-A": map[string][]int{
						"STEP-A": {0, 5, 0, 0, 0, 0, 0, 0},
						"STEP-B": {0, 0, 5, 0, 0, 0, 0, 0},
						"STEP-C": {0, 0, 0, 5, 0, 0, 0, 0},
					},
					"STEP-B": map[string][]int{
						"STEP-B": {0, 5, 0, 0, 0, 0, 0, 0},
						"STEP-C": {0, 0, 5, 0, 0, 0, 0, 0},
					},
					"STEP-C": map[string][]int{
						"STEP-C": {0, 5, 0, 0, 0, 0, 0, 0},
					},
				},
			},
		},
	}
	a := require.New(t)
	for _, test := range tests {
		partItems := make([]*PartItem, test.nItems)
		for i := 0; i < test.nItems; i++ {
			partItems[i] = NewPartItem(test.generateSerialNumber(i), test.stepOrdering)
			test.setStepsTime(test.stepOrdering, partItems[i])
		}
		items := make([]Item, test.nItems)
		for i := 0; i < test.nItems; i++ {
			items[i] = partItems[i]
		}
		master := NewDefaultMaster()
		out := master.Run(items, test.stepOrdering, test.nHouesInOneInterval, test.nInterval)
		a.Equal(test.want, out)
	}
}
