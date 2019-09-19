package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Reader struct {
	Feeder       chan FileContent
	basePath     string
	fileMetaChan chan Meta
	counterChan  chan bool
	readerCount  int
}

func NewReader(path string, readBufferSize int, readerCount int) Reader {
	if strings.HasSuffix(path, "/") {
		strings.TrimSuffix(path, "/")
	}
	feederChannel := make(chan FileContent, readBufferSize)
	readerFeeder := make(chan Meta, readerCount)
	return Reader{
		Feeder:       feederChannel,
		basePath:     path,
		fileMetaChan: readerFeeder,
		counterChan:  make(chan bool, readerCount),
		readerCount:  readerCount,
	}
}

func (r Reader) Start(operationMode string, currentIndex int64) {
	for i := 0; i < r.readerCount; i++ {
		go r.readAndFeed()
	}
	go r.timeTracker()
	r.feedFilesOfDirectory(operationMode, currentIndex)
}

func (r Reader) feedFilesOfDirectory(operationMode string, currentIndex int64) {
	var index int64 = 0
	err := filepath.Walk(r.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Callback for directory walk failed : "+path, err)
		}
		if !info.IsDir() {
			index++
			if operationMode == "restart" && index <= currentIndex {
				return nil
			}
			r.fileMetaChan <- Meta{
				FileInfo: info,
				Index:    index,
				Path:     strings.TrimSuffix(strings.TrimPrefix(path, r.basePath), info.Name()),
				FullPath: path,
			}
		}
		return nil
	})
	panicOnError(err, "error in listing file")
}

func panicOnError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}

func (r Reader) timeTracker() {
	count := 0
	timeInSeconds := 0
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-r.counterChan:
			count++
		case <-ticker.C:
			timeInSeconds++
			go func() { fmt.Print("\r", count, "/", timeInSeconds) }()
		}
	}
}

func (r Reader) readAndFeed() {
	for {
		info := <-r.fileMetaChan
		file, err := ioutil.ReadFile(info.FullPath)
		if err != nil {
			fmt.Println("Error in reading file with name "+info.Name()+"; Error : ", err)
		}
		content := NewFileContent(len(file), info.Name(), info.Index, info.Path, file)
		r.Feeder <- content
		r.counterChan <- true
	}
}
