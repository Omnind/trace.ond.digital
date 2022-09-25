package worker

import (
	"fmt"
	"time"
)

// Item is a struct that contains the information of a single item.
type Item interface {
	// GetSerialNumber returns the serial number of the item. The serial number should be unique.
	GetSerialNumber() string
	// GetTimeOfStep returns the time of the item.
	GetTimeOfStep(step string) time.Time
	// GetSteps returns the all steps.
	GetSteps() []string
}

type PartItem struct {
	// SerialNumber is the serial number of the item.
	serialNumber string
	// TimeOfStep is the time of the item.
	stepsOrdering []string
	// stepTime is the time of the item.
	stepTime map[string]time.Time
}

func NewPartItem(serialNumber string, stepsOrdering []string, times []time.Time) (*PartItem, error) {
	item := &PartItem{
		serialNumber:  serialNumber,
		stepsOrdering: stepsOrdering,
		stepTime:      make(map[string]time.Time),
	}

	if len(stepsOrdering) != len(times) {
		return nil, fmt.Errorf("the length of stepsOrdering and times should be the same")
	}

	for idx, step := range stepsOrdering {
		item.stepTime[step] = times[idx]
	}
	return item, nil
}

func (item *PartItem) GetSerialNumber() string {
	return item.serialNumber
}

func (item *PartItem) GetTimeOfStep(step string) time.Time {
	return item.stepTime[step]
}

func (item *PartItem) GetSteps() []string {
	return item.stepsOrdering
}
