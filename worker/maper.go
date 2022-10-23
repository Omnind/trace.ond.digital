package worker

import (
	"appledata/logger"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// worker is the member to run task in the map-reduce.
type worker struct {
	id int
	// Control panel.
	// errChan is the channel to notify the maper is done.
	errChan chan<- error
	wg      *sync.WaitGroup
}

func newWorker(id int, errChan chan<- error, wg *sync.WaitGroup) worker {
	return worker{
		id:      id,
		errChan: errChan,
		wg:      wg,
	}
}

type maper struct {
	worker
	// Data panel.
	data                 <-chan []Item
	middleResultColleter chan<- *ResultSet
	nHoursInOneInterval  int
	nInterval            int
}

// newMaper returns a new maper.
func newMaper(id int, errChan chan<- error, wg *sync.WaitGroup, data <-chan []Item, middleResultColleter chan<- *ResultSet, nHoursInOneInterval, nInterval int) *maper {
	return &maper{
		worker:               newWorker(id, errChan, wg),
		data:                 data,
		middleResultColleter: middleResultColleter,
		nHoursInOneInterval:  nHoursInOneInterval,
		nInterval:            nInterval,
	}
}

func (m *maper) run(stepOrdering []string) {
	defer func() {
		if r := recover(); r != nil {
			err := errors.Errorf("panic in maper: %v", r)
			m.errChan <- err
		}
		logger.Info("maper stopped", zap.Int("maper", m.id))
		m.wg.Done()
	}()
	logger.Info("maper started", zap.Int("maper", m.id))
	for {
		items, ok := <-m.data
		if !ok {
			close(m.middleResultColleter)
			return
		}
		logger.Debug("maper received data", zap.Int("maper", m.id), zap.Int("nItem", len(items)))
		middleResults := handleItems(m.nInterval, m.nHoursInOneInterval, stepOrdering, items)
		// NOTE: The middleResultColleter is a buffered channel, so it will not block. But send data to a closed channel will panic.
		m.middleResultColleter <- middleResults
	}
}

func handleItems(nInterval int, nHouesInOneInterval int, stepsOrdering []string, items []Item) *ResultSet {
	mr := &ResultSet{
		PassStepsTimeNumber: make(map[string]ToStepIntervalSpent),
		FailStepsTimeNumber: make(StepsTimeCalculator),
	}
	if len(items) == 0 {
		return nil
	}
	mr.PassStepsTimeNumber.initByStepsOrdering(stepsOrdering, nInterval)
	mr.FailStepsTimeNumber.initByStepsOrdering(stepsOrdering, nInterval)

	for _, item := range items {
		nStep := len(stepsOrdering)
		for i := 0; i < nStep; i++ {
			for j := i; j < nStep; j++ {
				fromStep := stepsOrdering[i]
				toStep := stepsOrdering[j]
				fromStepInfo, ok := item.GetStep(fromStep)
				if !ok {
					// The item may lack some steps.
					continue
				}
				toStepInfo, ok := item.GetStep(toStep)
				if !ok {
					// The item may lack some steps.
					continue
				}
				beginTime := fromStepInfo.GetBeginTime()
				endTime := toStepInfo.GetEndTime()
				tmp := endTime.Sub(beginTime)
				timeSpentInterval := int(tmp.Seconds() / float64(nHouesInOneInterval*3600 /* to second*/))
				timeSpentInterval = min(timeSpentInterval, nInterval)
				timeSpentInterval = max(timeSpentInterval, 0)
				if fromStepInfo.GetStatus() == StepPass && toStepInfo.GetStatus() == StepPass {
					mr.PassStepsTimeNumber[fromStep][toStep][timeSpentInterval+1] += 1
				} else {
					mr.FailStepsTimeNumber[fromStep][toStep][timeSpentInterval+1] += 1
				}
			}
		}
	}
	return mr
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
