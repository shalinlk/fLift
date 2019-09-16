package file

import (
	"fmt"
	"io/ioutil"
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
		//basePath = fmt.Sprintf("%s/", basePath)
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
	r.feedFilesOfDirectory(0, operationMode, currentIndex, "/")
}

func (r Reader) feedFilesOfDirectory(index int64, operationMode string, currentIndex int64, path string) int64 {
	files, err := ioutil.ReadDir(r.basePath + path)
	panicOnError(err, "error in listing file")
	for _, file := range files {
		if file.IsDir() {
			index = r.feedFilesOfDirectory(index, operationMode, currentIndex, fmt.Sprintf("%s%s/", path, file.Name()))
		} else {
			index++
			if operationMode == "restart" && index <= currentIndex {
				continue
			}
			r.fileMetaChan <- Meta{
				FileInfo: file,
				Index:    index,
				Path:     path,
			}
		}
	}
	return index
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
		file, err := ioutil.ReadFile(r.basePath + info.Name())
		if err != nil {
			fmt.Println("Error in reading file with name "+info.Name()+"; Error : ", err)
		}
		content := NewFileContent(len(file), info.Name(), info.Index, info.Path, file)
		r.Feeder <- content
		r.counterChan <- true
	}
}
