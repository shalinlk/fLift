package file

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Writer struct {
	basePath string
}

func NewWriter(basePath string) Writer {
	basePath = strings.TrimSpace(basePath)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	return Writer{basePath}
}
func (w *Writer) writeToFile(content FileContent) {
	err := ioutil.WriteFile(w.basePath+content.Name, content.getBytes(), 0644)
	if err != nil {
		fmt.Print("Error in writing file : " + content.Name)
	}
}
