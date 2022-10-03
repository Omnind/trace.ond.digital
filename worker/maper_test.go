package worker

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandleItems(t *testing.T) {
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
			nItems:               1,
			nInterval:            7,
			nHouesInOneInterval:  24,
			generateSerialNumber: commonGenerateSerialNumber,
			setStepsTime:         commonSetStepsTime,
			want: &ResultSet{
				StepsTimeNumber: map[string]ToStepIntervalSpent{
					"STEP-A": map[string][]int{
						"STEP-A": {0, 1, 0, 0, 0, 0, 0, 0},
						"STEP-B": {0, 0, 1, 0, 0, 0, 0, 0},
						"STEP-C": {0, 0, 0, 1, 0, 0, 0, 0},
					},
					"STEP-B": map[string][]int{
						"STEP-B": {0, 1, 0, 0, 0, 0, 0, 0},
						"STEP-C": {0, 0, 1, 0, 0, 0, 0, 0},
					},
					"STEP-C": map[string][]int{
						"STEP-C": {0, 1, 0, 0, 0, 0, 0, 0},
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
		out := handleItems(test.nInterval, test.nHouesInOneInterval, test.stepOrdering, items)
		a.Equal(test.want, out)
	}
}

// commonGenerateSerialNumber is a helper function to generate serial number.
// It use `i` as serial number.
func commonGenerateSerialNumber(i int) string {
	return strconv.Itoa(i)
}

// commonSetStepsTime is a helper function to generate step time.
// Each step will use 24 hours.
func commonSetStepsTime(stepsOrdering []string, item *PartItem) {
	baseTime := time.Now()
	beginTime := baseTime
	for _, step := range stepsOrdering {
		endTime := beginTime.Add(time.Duration(24) * time.Hour)
		item.SetStep(step, NewStep(step, beginTime, endTime, StepSuccess))
		beginTime = endTime
	}
}
