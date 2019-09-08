package file

import (
	"fmt"
	"io/ioutil"
	"strings"
	. "time"
)

type Writer struct {
	basePath     string
	count        int
	consumerChan chan FileContent
}

func NewWriter(basePath string, consumerChan chan FileContent) Writer {
	basePath = strings.TrimSpace(basePath)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	return Writer{
		basePath:     basePath,
		count:        0,
		consumerChan: consumerChan,
	}
}

func (w Writer) ConsumeAndWrite() {
	ticker := NewTicker(Second)
	for {
		select {
		case content := <-w.consumerChan:
			w.WriteToFile(content)
		case t := <-ticker.C:
			fmt.Println("Written : ", w.count, " by ", t)
		}
	}
}
func (w *Writer) WriteToFile(content FileContent) {
	err := ioutil.WriteFile(w.basePath+content.Name, content.getBytes(), 0644)
	if err != nil {
		fmt.Print("Error in writing file : " + content.Name)
	} else {
		w.count++
	}
}