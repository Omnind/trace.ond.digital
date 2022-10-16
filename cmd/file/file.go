package file

import (
	"appledata/logger"
	"appledata/worker"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func ReadFileAndConvertToItem(path string, stepOrdering []string, stepShortName, fullStepName string, result chan<- map[string]worker.Item, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(path)
	if err != nil {
		logger.Error("failed to open file", zap.Error(err))
		errChan <- errors.Wrapf(err, "failed to open file %s", path)
		return
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ','
	// We should find the serialId, startTime and endTime column idx.
	columnNames, err := csvReader.Read()
	if err != nil {
		logger.Error("failed to read column names", zap.Error(err))
		errChan <- errors.Wrapf(err, "failed to read the first line of file %s", path)
		return
	}

	serialNumberColumnName := getSerialNumberColumnName()
	beginTimeColumnName := getTimeColumnName(stepShortName, true)
	endTimeColumnName := getTimeColumnName(stepShortName, false)
	statusColumnName := getStatusColumnName(stepShortName)
	var serialNumberIdx int
	var beginTimeIdx int
	var endTimeIdx int
	var statusIdx int
	for i, columnName := range columnNames {
		// We need to remove the BOM.
		// https://zasy.github.io/2018/09/28/tx-06/
		columnName = strings.Replace(columnName, "\ufeff", "", -1)
		if columnName == serialNumberColumnName {
			serialNumberIdx = i
		} else if columnName == beginTimeColumnName {
			beginTimeIdx = i
		} else if columnName == endTimeColumnName {
			endTimeIdx = i
		} else if columnName == statusColumnName {
			statusIdx = i
		} else {
			logger.Info("not handle column", zap.String("column name", columnName))
		}
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Error("failed to read all records", zap.Error(err))
		errChan <- errors.Wrapf(err, "failed to read all records of file %s", path)
		return
	}

	itemMap := make(map[string]worker.Item)
	for _, record := range records {
		serialNumber := record[serialNumberIdx]
		beginTimeStr := record[beginTimeIdx]
		beginTime, err := normalizeTime(beginTimeStr)
		if err != nil {
			logger.Error("failed to parse begin time", zap.String("beginTimeStr", beginTimeStr), zap.Error(err))
			errChan <- errors.Wrapf(err, "failed to parse begin time %s", beginTimeStr)
			return
		}
		endTimeStr := record[endTimeIdx]
		endTime, err := normalizeTime(endTimeStr)
		if err != nil {
			logger.Error("failed to parse end time %s", zap.String("endTimeStr", endTimeStr))
			errChan <- errors.Wrapf(err, "failed to parse end time %s", endTimeStr)
			return
		}
		// TODO(zp): handle the status
		statusStr := record[statusIdx]
		status := worker.StepFail
		if statusStr == "passed" {
			status = worker.StepPass
		}
		item := worker.NewPartItem(serialNumber, stepOrdering)
		stepInfo := worker.NewStep(fullStepName, beginTime, endTime, status)
		item.SetStep(fullStepName, stepInfo)
		itemMap[serialNumber] = item
	}
	result <- itemMap
}

func getTimeColumnName(shortStepName string, begin bool) string {
	if len(shortStepName) == 0 {
		return ""
	}
	var specs []string
	specs = append(specs, shortStepName)
	specs = append(specs, "insight")
	specs = append(specs, "test_attributes")
	if begin {
		specs = append(specs, "uut_start")
	} else {
		specs = append(specs, "uut_stop")
	}
	return strings.Join(specs, ".")
}

func getSerialNumberColumnName() string {
	return "root_serial"
}

func getStatusColumnName(shortStepName string) string {
	return "result"
}

func normalizeTime(timeStr string) (time.Time, error) {
	specs := strings.Split(timeStr, " ")
	date := specs[0]
	eles := strings.Split(date, "/")
	if hours, err := strconv.Atoi(eles[1]); err != nil {
		return time.Time{}, err
	} else {
		if hours < 10 {
			eles[1] = "0" + eles[1]
		}
	}
	if miniute, err := strconv.Atoi(eles[2]); err != nil {
		return time.Time{}, err
	} else {
		if miniute < 10 {
			eles[2] = "0" + eles[2]
		}
	}
	correct := strings.Join(eles, "/") + " " + specs[1]
	return time.Parse("2006/01/02 15:04", correct)
}
