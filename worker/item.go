package worker

import (
	"time"
)

// Item is a struct that contains the information of a single item.
type Item interface {
	// GetSerialNumber returns the serial number of the item. The serial number should be unique.
	GetSerialNumber() string
	// GetStep returns the step info of the specify setp.
	GetStep(stepName string) (*Step, bool)
	// GetAllSteps returns the steps info of the item.
	GetAllSteps() []*Step
	// GetStepsOrdering returns the all steps by order.
	GetStepsOrdering() []string
}

// stepStatus is the status of the step
type stepStatus string

var (
	StepSuccess stepStatus = "SUCCESS"
	StepFail    stepStatus = "FAIL"
)

// Step includes the infomation of one step.
type Step struct {
	name      string
	beginTime time.Time
	endTime   time.Time
	status    stepStatus
}

func (s *Step) GetName() string {
	return s.name
}

func (s *Step) GetBeginTime() time.Time {
	return s.beginTime
}

func (s *Step) GetEndTime() time.Time {
	return s.endTime
}

func (s *Step) GetStatus() stepStatus {
	return s.status
}

// NewStep returns a new step.
func NewStep(name string, beginTime time.Time, endTime time.Time, status stepStatus) *Step {
	return &Step{
		name:      name,
		beginTime: beginTime,
		endTime:   endTime,
		status:    status,
	}
}

type PartItem struct {
	// SerialNumber is the serial number of the item.
	serialNumber string
	// TimeOfStep is the time of the item.
	stepsOrdering []string
	// steps is the map of the each step.
	steps map[string]*Step
}

func NewPartItem(serialNumber string, stepsOrdering []string) *PartItem {
	return &PartItem{
		serialNumber:  serialNumber,
		stepsOrdering: stepsOrdering,
		steps:         make(map[string]*Step),
	}
}

func (item *PartItem) SetStep(stepName string, step *Step) {
	item.steps[stepName] = step
}

func (item *PartItem) GetSerialNumber() string {
	return item.serialNumber
}

func (item *PartItem) GetStep(stepName string) (*Step, bool) {
	step, ok := item.steps[stepName]
	return step, ok
}

func (item *PartItem) GetStepsOrdering() []string {
	return item.stepsOrdering
}

func (item *PartItem) GetAllSteps() []*Step {
	steps := make([]*Step, 0, len(item.steps))
	for _, step := range item.steps {
		steps = append(steps, step)
	}
	return steps
}
