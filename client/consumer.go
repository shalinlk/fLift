package client

import (
	"errors"
	"fmt"

	"github.com/shalinlk/fLift/file"
)

func connectionFactory(connType, host string) (Client, error) {
	if connType == "tcp" {
		return NewTCPClient(host), nil
	}
	return nil, errors.New("connection type not defined")
}
func StartConsumer(host, baseFilePath, connType string, writeBufferSize int) {
	socket, err := connectionFactory(host, connType)
	checkAndPanicOnError(err)
	consumerChannel := make(chan file.FileContent, writeBufferSize)
	socket.Start(consumerChannel)
}

func checkAndPanicOnError(err error) {
	if err != nil {
		fmt.Println("Error : ", err)
		panic(err)
	}
}
