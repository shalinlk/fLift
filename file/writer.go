package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	. "time"
)

type Writer struct {
	basePath         string
	consumerChan     chan FileContent
	agentCount       int
	counterChan      chan bool
	directoryMap     map[string]bool
	directoryMapLock *sync.Mutex
}

func NewWriter(basePath string, consumerChan chan FileContent, agentCount int) Writer {
	basePath = strings.TrimSpace(basePath)
	if strings.HasSuffix(basePath, "/") {
		strings.TrimSuffix(basePath, "/")
	}
	return Writer{
		basePath:         basePath,
		consumerChan:     consumerChan,
		agentCount:       agentCount,
		counterChan:      make(chan bool, agentCount*2),
		directoryMap:     make(map[string]bool),
		directoryMapLock: &sync.Mutex{},
	}
}

func (w Writer) StartWriters() {
	for i := 0; i < w.agentCount; i++ {
		go func() {
			for {
				w.WriteToFile(<-w.consumerChan)
			}
		}()
	}
	go w.timeTracker()
}

func (w *Writer) WriteToFile(content FileContent) {
	w.createDirectoryIfDoesNotExist(content.Path)
	err := ioutil.WriteFile(w.basePath+content.Path+content.Name, content.getBytes(), 0644)
	if err != nil {
		fmt.Print("Error in writing file : " + content.Name)
	} else {
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

func (w Writer) createDirectoryIfDoesNotExist(path string) {
	if path == "/" {
		return
	}
	w.directoryMapLock.Lock()
	_, exist := w.directoryMap[path]
	w.directoryMapLock.Unlock()
	if exist {
		return
	}
	_ = os.MkdirAll(w.basePath+path, 0777)
	w.directoryMapLock.Lock()
	w.directoryMap[path] = true
	w.directoryMapLock.Unlock()
}
