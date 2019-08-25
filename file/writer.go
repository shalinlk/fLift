package file

import (
	"io/ioutil"
	"strings"
)

type writer struct {
	basePath string
}

func NewWriter(basePath string) writer {
	basePath = strings.TrimSpace(basePath)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	return writer{basePath}
}
func (w *writer) writeToFile(content FileContent) {
	ioutil.WriteFile(w.basePath+content.Name, content.getBytes(), 0644)
}
