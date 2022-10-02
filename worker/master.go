package worker

import (
	"appledata/logger"
	"sync"
	"time"
)

const (
	MapWorkerNumber       = 3
	ChannelBufferSize     = 5
	BatchSize             = 1024
	CollectResultInterval = 300 * time.Millisecond
)

// resultSet is the result of the collector.
type resultSet struct {
	// stepsTimeNumber is the number of time interval of steps.
	stepsTimeNumber stepsTimeCalculator
}

type toStepIntervalSpent map[string][]int

type stepsTimeCalculator map[string]toStepIntervalSpent

// initByStepsOrdering will init the stepsTimeNumber by the steps ordering.
func (s stepsTimeCalculator) initByStepsOrdering(stepsOrdering []string, nInterval int) {
	for i := 0; i < len(stepsOrdering); i++ {
		fromStep := stepsOrdering[i]
		s[fromStep] = make(toStepIntervalSpent)
		for j := i; j < len(stepsOrdering); j++ {
			toStep := stepsOrdering[j]
			// We make the slice with length `nInterval+1` because we should record the times that item used more than interval * nInterval.
			// And the capactiy of slice is the same as the length.
			s[fromStep][toStep] = make([]int, nInterval+1)
		}
	}
}

type Master struct {
	nMapWorker    int
	errChan       []chan error
	data          []chan []Item
	collect       []chan *resultSet
	wg            *sync.WaitGroup
	workerAlive   []bool
	aliveWorker   int
	result        *resultSet
	collectorDone chan struct{}
	finished      chan struct{}
}

func NewDefaultMaster() *Master {
	return &Master{
		nMapWorker:  MapWorkerNumber,
		aliveWorker: MapWorkerNumber,
		result: &resultSet{
			stepsTimeNumber: make(stepsTimeCalculator),
		},
	}
}

func (m *Master) Run(items []Item, stepOrdering []string, nHoursInOneInterval, nInterval int) *resultSet {
	if len(items) == 0 {
		return nil
	}
	// We should init the result before we start the worker.
	m.result.stepsTimeNumber.initByStepsOrdering(stepOrdering, nInterval)

	// We prepare the worker first.
	m.errChan = make([]chan error, m.nMapWorker)
	m.data = make([]chan []Item, m.nMapWorker)
	m.collect = make([]chan *resultSet, m.nMapWorker)
	m.wg = &sync.WaitGroup{}
	m.workerAlive = make([]bool, m.nMapWorker)
	m.collectorDone = make(chan struct{})
	// We should make the finished channel with buffer size 1 so that the colletor can return quickly.
	m.finished = make(chan struct{}, 1)

	// Start the collector.
	logger.Info("master start the collector")
	go m.runColletor()
	defer func() {
		m.collectorDone <- struct{}{}
		<-m.finished
	}()

	for i := 0; i < m.nMapWorker; i++ {
		m.errChan[i] = make(chan error, 1)
		m.data[i] = make(chan []Item, ChannelBufferSize)
		m.collect[i] = make(chan *resultSet, ChannelBufferSize)
		m.workerAlive[i] = true
		maper := newMaper(i, m.errChan[i], m.wg, m.data[i], m.collect[i], nHoursInOneInterval, nInterval)
		m.wg.Add(1)
		go maper.run(stepOrdering)
		time.Sleep(100 * time.Millisecond)
	}

	last := 0
	for {
		for i := 0; i < m.nMapWorker; i++ {
			select {
			case <-m.errChan[i]:
				m.aliveWorker--
				m.workerAlive[i] = false
			default:
			}
		}
		if m.aliveWorker == 0 {
			// all worker is dead, we should stop the master.
			logger.Info("all workers had shutdown, we should stop the master.")
			return m.result
		}
		if last >= len(items) {
			// close all data channel to notify the worker to stop.
			for i := 0; i < m.nMapWorker; i++ {
				close(m.data[i])
			}
			return m.result
		}
		// We handle BatchSize items each time.
		batchData := make([][]Item, m.nMapWorker)
		for i := 0; i < m.nMapWorker; i++ {
			batchData[i] = make([]Item, 0, BatchSize)
		}
		for i := last; i < last+BatchSize && i < len(items); i++ {
			item := items[i]
			workerIdx := m.chooseWorker(item)
			batchData[workerIdx] = append(batchData[workerIdx], item)
		}
		// Send to worker
		for i := 0; i < m.nMapWorker; i++ {
			if m.workerAlive[i] {
				m.data[i] <- batchData[i]
			}
		}
		last += BatchSize
	}
}

func (m *Master) chooseWorker(item Item) int {
	bs := item.GetSerialNumber()
	originalWorker := int(bs[len(bs)-1]) % m.nMapWorker
	if !m.workerAlive[originalWorker] {
		idx := 0
		for step := 0; step < originalWorker; step++ {
			idx = (idx + 1) % m.nMapWorker
			for !m.workerAlive[idx] {
				idx = (idx + 1) % m.nMapWorker
			}
		}
		return idx
	}
	return originalWorker
}

func (m *Master) Wait() {
	// All the data channel should be closed, we will wait all data.
	for i := 0; i < m.nMapWorker; i++ {
		// For-range will read all data from the channel until the channel is closed.
		for middleData := range m.collect[i] {
			if middleData == nil {
				continue
			}
			for fromStep, toStepTimeSpent := range middleData.stepsTimeNumber {
				for toStep, number := range toStepTimeSpent {
					for i, v := range number {
						m.result.stepsTimeNumber[fromStep][toStep][i] += v
					}
				}
			}
		}
	}
	logger.Info("master waiting all workers to shutdown...")
	m.wg.Wait()
}

func (m *Master) runColletor() {
	for {
		select {
		case <-m.collectorDone:
			// All the data channel should be closed, we will wait all data.
			for i := 0; i < m.nMapWorker; i++ {
				// For-range will read all data from the channel until the channel is closed.
				for middleData := range m.collect[i] {
					if middleData == nil {
						continue
					}
					for fromStep, toStepTimeSpent := range middleData.stepsTimeNumber {
						for toStep, number := range toStepTimeSpent {
							for i, v := range number {
								m.result.stepsTimeNumber[fromStep][toStep][i] += v
							}
						}
					}
				}
			}
			logger.Info("colltor waiting all workers to shutdown...")
			m.wg.Wait()
			m.finished <- struct{}{}
			return
		default:
			for i := 0; i < m.nMapWorker; i++ {
				select {
				case middleData := <-m.collect[i]:
					if middleData == nil {
						continue
					}
					for fromStep, toStepTimeSpent := range middleData.stepsTimeNumber {
						for toStep, number := range toStepTimeSpent {
							for i, v := range number {
								m.result.stepsTimeNumber[fromStep][toStep][i] += v
							}
						}
					}
				default:
				}
			}
		}
	}
}
