package file

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Writer struct {
	basePath string
	count int
}

func NewWriter(basePath string) Writer {
	basePath = strings.TrimSpace(basePath)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	return Writer{basePath, 0}
}
func (w *Writer) WriteToFile(content FileContent) {
	err := ioutil.WriteFile(w.basePath+content.Name, content.getBytes(), 0644)
	if err != nil {
		fmt.Print("Error in writing file : " + content.Name)
	}else{
		w.count++
		fmt.Println("\rWritten : ", w.count)
	}
}
