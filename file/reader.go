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
	readerFeeder chan os.FileInfo
	counterChan  chan bool
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
		readerFeeder: readerFeeder,
		counterChan:  make(chan bool, readerCount),
	}
}

func (r Reader) Start() {
	for i := 0; i < 10; i++ {
		go r.readAndFeed2()
	}
	go r.tracker()
	files, err := ioutil.ReadDir(r.path)
	panicOnError(err, "error in listing file")
	fmt.Println("Reading")
	for _, file := range files {
		r.readerFeeder <- file
	}
}

func panicOnError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}

func (r Reader) readAndFeed2() {
	for { r.readAndFeed(<- r.readerFeeder)}
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

func (r Reader) readAndFeed(info os.FileInfo) {
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
