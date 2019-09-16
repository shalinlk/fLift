package client

import (
	"errors"
	"fmt"

	. "github.com/shalinlk/fLift/file"
)

func connectionFactory(connType, host string, consumerChan chan FileContent, concurrentConnections int) (Client, error) {
	if connType == "tcp" {
		return NewTCPClient(host, consumerChan, concurrentConnections), nil
	}
	return nil, errors.New("connection type not defined")
}
func StartConsumer(host, writeFilePath, connType string, writeBufferSize int, agentCount int, concurrentConnections int) {
	if concurrentConnections < 1{
		concurrentConnections = 1
	}
	consumerChannel := make(chan FileContent, writeBufferSize)
	socket, err := connectionFactory(connType, host, consumerChannel, concurrentConnections)
	checkAndPanicOnError(err)
	writer := NewWriter(writeFilePath, consumerChannel, agentCount) //todo : writer should have been a dependency of client
	writer.StartWriters()
	socket.Start()
}

func checkAndPanicOnError(err error) {
	if err != nil {
		fmt.Println("Error : ", err)
		panic(err)
	}
}
