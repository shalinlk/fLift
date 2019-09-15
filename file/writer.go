package file

import (
	"fmt"
	"io/ioutil"
	"strings"
	. "time"
)

type Writer struct {
	basePath     string
	consumerChan chan FileContent
	agentCount   int
	counterChan  chan bool
}

func NewWriter(basePath string, consumerChan chan FileContent, agentCount int) Writer {
	basePath = strings.TrimSpace(basePath)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	return Writer{
		basePath:     basePath,
		consumerChan: consumerChan,
		agentCount:   agentCount,
		counterChan:  make(chan bool, agentCount*2),
	}
}

func (w Writer) StartWriters() {
	for i := 0; i < w.agentCount; i++ {
		go w.startWriter()
	}
	go w.timeTracker()
}

func (w Writer) startWriter() {
	for {
		w.WriteToFile(<-w.consumerChan)
	}
}

func (w *Writer) WriteToFile(content FileContent) {
	err := ioutil.WriteFile(w.basePath+content.Name, content.getBytes(), 0644)
	if err != nil {
		fmt.Print("Error in writing file : " + content.Name)
	}else {
		w.counterChan <- true
	}
}

func (w Writer) timeTracker() {
	count := 0
	timeInSeconds := 0
	ticker := NewTicker(Second * 1)
	for {
		select {
		case <-w.counterChan:
			count++
		case <-ticker.C:
			timeInSeconds++
			go func() { fmt.Print("\r", count, "/", timeInSeconds) }()
		}
	}
}
