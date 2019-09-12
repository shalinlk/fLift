package client

import (
	"errors"
	"fmt"

	. "github.com/shalinlk/fLift/file"
)

func connectionFactory(connType, host string, consumerChan chan FileContent) (Client, error) {
	if connType == "tcp" {
		return NewTCPClient(host, consumerChan), nil
	}
	return nil, errors.New("connection type not defined")
}
func StartConsumer(host, writeFilePath, connType string, writeBufferSize int, agentCount int) {
	consumerChannel := make(chan FileContent, writeBufferSize)
	socket, err := connectionFactory(connType, host, consumerChannel)
	checkAndPanicOnError(err)
	writer := NewWriter(writeFilePath, consumerChannel, agentCount)
	writer.StartWriters()
	socket.Start()
}

func checkAndPanicOnError(err error) {
	if err != nil {
		fmt.Println("Error : ", err)
		panic(err)
	}
}
