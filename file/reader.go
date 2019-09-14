package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Reader struct {
	Feeder       chan FileContent
	path         string
	fileMetaChan chan os.FileInfo
	counterChan  chan bool
	readerCount  int
}

func NewReader(path string, readBufferSize int, readerCount int) Reader {
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%s/", path)
	}
	feederChannel := make(chan FileContent, readBufferSize)
	readerFeeder := make(chan os.FileInfo, readerCount)
	return Reader{
		Feeder:       feederChannel,
		path:         path,
		fileMetaChan: readerFeeder,
		counterChan:  make(chan bool, readerCount),
		readerCount:  readerCount,
	}
}

func (r Reader) Start() {
	for i := 0; i < r.readerCount; i++ {
		go r.readAndFeed()
	}
	go r.tracker()
	files, err := ioutil.ReadDir(r.path) // os commands
	panicOnError(err, "error in listing file")
	fmt.Println("Reading")
	for _, file := range files {
		r.fileMetaChan <- file
	}
}

func panicOnError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}

func (r Reader) tracker() {
	count := 0
	second := 0
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-r.counterChan:
			count++
		case <-ticker.C:
			second++
			fmt.Println(count, "/", second)
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
		content := FileContent{
			Size:    len(file),
			Name:    info.Name(),
			Content: file,
		}
		r.Feeder <- content
		r.counterChan <- true
	}
}
