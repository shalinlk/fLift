package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Reader struct {
	Feeder       chan FileContent
	basePath     string
	fileMetaChan chan Meta
	counterChan  chan bool
	readerCount  int
	batchSize    int
	stopChan     chan bool
}

func NewReader(path string, readBufferSize int, readerCount int, batchSize int) Reader {
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
		batchSize:    batchSize,
		stopChan:     make(chan bool),
	}
}

func (r Reader) Start(operationMode string, currentIndex int64) {
	var wg sync.WaitGroup
	wg.Add(r.readerCount)
	for i := 0; i < r.readerCount; i++ {
		go r.readAndFeed(&wg)
	}
	go r.timeTracker()
	var index int64 = 0
	r.feedFilesOfDirectory(operationMode, currentIndex, index, "/")
	close(r.fileMetaChan)
	wg.Wait()
	r.stopChan <- true
	close(r.counterChan)
	close(r.Feeder)
	fmt.Println("Reading files finished at : ", time.Now().UnixNano()/int64(time.Millisecond))
}

func (r *Reader) feedFilesOfDirectory(operationMode string, currentIndex int64, index int64, path string) int64 {

	directory, directoryError := os.Open(r.basePath + path)
	if directoryError != nil {
		fmt.Println("Error in opening directory ", directoryError)
	}
	for {
		metaOfFiles, err := directory.Readdir(r.batchSize)
		if err != nil {
			fmt.Println("Truncating Meta Reader, time : ", time.Now().UnixNano()/int64(time.Millisecond))
			if err != io.EOF {
				fmt.Println("Truncating Meta reading because of error : ", err)
			}
			break
		}
		for _, info := range metaOfFiles {
			if ! info.IsDir() {
				index++
				if operationMode == "restart" && index <= currentIndex {
					continue
				}
				r.fileMetaChan <- Meta{
					FileInfo: info,
					Index:    index,
					Path:     path,
				}
			}
		}
	}
	return index
}

func (r *Reader) timeTracker() {
	count := 0
	timeInSeconds := 0
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-r.counterChan:
			count++
		case <-ticker.C:
			timeInSeconds++
			go func(timeSec int) { fmt.Print("\r", count, "/", timeSec) }(timeInSeconds)
		case <-r.stopChan:
			ticker.Stop()
			return
		}
	}
}

func (r Reader) readAndFeed(wg *sync.WaitGroup) {
	for info := range r.fileMetaChan {
		file, err := ioutil.ReadFile(r.basePath + info.Path + info.FileInfo.Name())
		if err != nil {
			fmt.Println("Error in reading file with name "+info.Name()+"; Error : ", err)
		}
		content := NewFileContent(len(file), info.Name(), info.Index, info.Path, file)
		r.Feeder <- content
		r.counterChan <- true
	}
	wg.Done()
}
