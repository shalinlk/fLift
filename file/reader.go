package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Reader struct {
	Feeder chan FileContent
	path   string
}

func NewReader(path string, readBufferSize int) Reader {
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%s/", path)
	}
	feederChannel := make(chan FileContent, readBufferSize)
	return Reader{Feeder: feederChannel, path:path}
}

func (r Reader) Start()  {
	files, err := ioutil.ReadDir(r.path)
	panicOnError(err, "error in listing file")
	count := 0;
	fmt.Println("Reading...")
	for _,file:= range files   {
		r.readAndFeed(file)
		count++
		fmt.Println("\rNumber of files Read : ", count)
	}
}

func panicOnError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}

func (r Reader) readAndFeed(info os.FileInfo) {
	file, err := ioutil.ReadFile(r.path + info.Name())
	if err != nil {
		fmt.Println("Error in reading file with name " + info.Name() + "; Error : ", err)
	}
	content := FileContent{
		Size:    len(file),
		Name:    info.Name(),
		Content: file,
	}
	r.Feeder <- content
}