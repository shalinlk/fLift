package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type StatusTracker struct {
	ticker            *time.Ticker
	index             int64
	statusTrackerChan chan int64
	fileHandle        *os.File
	keepStatus        bool
}

func NewStatusTracker(maxClients int, flushInterval int, operationMode string, keepStatus bool) StatusTracker {
	const ModeRestart = "restart"
	const ModeStart = "start"
	statusTrackerFileName := "status.txt"
	var index int64 = 0
	var ticker *time.Ticker = nil
	var fileHandle *os.File = nil

	if operationMode != ModeStart && operationMode != ModeRestart {
		panic("operationMode should be either " + ModeStart + " / " + ModeRestart)
	}
	if keepStatus {
		if operationMode == ModeStart {
			_ = os.Remove(statusTrackerFileName)
		}
		fileHandle, _ = os.OpenFile(statusTrackerFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if operationMode == ModeRestart {
			contentInBytes, err := ioutil.ReadFile(statusTrackerFileName)
			if err == nil {
				index, err = strconv.ParseInt(string(contentInBytes), 10, 64)
				if err != nil {
					index = 0
				}
			}
		}
		ticker = time.NewTicker(time.Second * time.Duration(flushInterval))
	}

	return StatusTracker{
		ticker:            ticker,
		index:             index,
		statusTrackerChan: make(chan int64, maxClients),
		fileHandle:        fileHandle,
	}
}

func (s *StatusTracker) CurrentIndex() int64 {
	return s.index
}

func (s *StatusTracker) StatusTrackerChan() chan int64 {
	return s.statusTrackerChan
}

func (s *StatusTracker) Start() {
	if s.keepStatus {
		go s.trackStatus()
	} else {
		go s.spoilStatus()
	}
}

func (s *StatusTracker) spoilStatus() {
	for {
		<-s.statusTrackerChan
	}
}

func (s *StatusTracker) trackStatus() {
	for {
		select {
		case index := <-s.statusTrackerChan:
			if index > s.index {
				s.index = index
			}
		case <-s.ticker.C:
			_, err := s.fileHandle.WriteAt([]byte(strconv.FormatInt(s.index, 10)), 0)
			if err != nil {
				fmt.Println("flushing status to file tracker failed ", err)
			}
		}
	}
}
