package file

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Reader struct {
	Feeder        chan FileContent
	path          string
	fileMetaChan  chan Meta
	counterChan   chan bool
	readerCount   int
}

func NewReader(path string, readBufferSize int, readerCount int ) Reader {
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%s/", path)
	}
	feederChannel := make(chan FileContent, readBufferSize)
	readerFeeder := make(chan Meta, readerCount)
	return Reader{
		Feeder:       feederChannel,
		path:         path,
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
	files, err := ioutil.ReadDir(r.path) // os commands
	panicOnError(err, "error in listing file")
	fmt.Println("Reading")
	var index int64 = 0
	for _, file := range files {
		index++
		if operationMode  == "restart" && index <=  currentIndex{
			continue
		}
		r.fileMetaChan <- Meta{
			FileInfo: file,
			Index:    index,
			Path:     file.Name(),
		}
	}
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
		file, err := ioutil.ReadFile(r.path + info.Name())
		if err != nil {
			fmt.Println("Error in reading file with name "+info.Name()+"; Error : ", err)
		}
		content := NewFileContent(len(file), info.Name(), info.Index, info.Path, file)
		r.Feeder <- content
		r.counterChan <- true
	}
}
