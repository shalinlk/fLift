package client

import (
	"errors"
	"fmt"

	. "github.com/shalinlk/fLift/file"
)

func connectionFactory(connType, host string) (Client, error) {
	if connType == "tcp" {
		return NewTCPClient(host), nil
	}
	return nil, errors.New("connection type not defined")
}
func StartConsumer(host, writeFilePath, connType string, writeBufferSize int) {
	socket, err := connectionFactory(connType, host)
	checkAndPanicOnError(err)
	consumerChannel := make(chan FileContent, writeBufferSize)
	writer := NewWriter(writeFilePath)
	go consumeAndWrite(consumerChannel, writer)
	socket.Start(consumerChannel)
}

func consumeAndWrite(contents chan FileContent, writer Writer) {
	for {
		select {
		case content := <- contents:
			writer.WriteToFile(content)
		}
	}
}

func checkAndPanicOnError(err error) {
	if err != nil {
		fmt.Println("Error : ", err)
		panic(err)
	}
}
