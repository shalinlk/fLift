package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/shalinlk/fLift/file"
	"github.com/shalinlk/fLift/utils"
)

type TCPClient struct {
	clientId string
	socket   net.Conn
}

func NewTCPClient(serverHost string) TCPClient {
	connection, err := net.Dial("tcp", serverHost)
	if err != nil {
		fmt.Println("connection failed to " + serverHost)
		panic(err)
	}
	return TCPClient{socket: connection}
}

func (c TCPClient) Register() {
	c.socket.Write([]byte(utils.FillUpForCommad("REGISTER")))
	clientIdBuffer := make([]byte, 60)
	c.socket.Read(clientIdBuffer)
	c.clientId = strings.Trim(string(clientIdBuffer), ":")
}

func (c TCPClient) Start(consumer chan<- file.FileContent) {
	c.socket.Write([]byte(utils.FillUpForCommad("START")))
	c.readAndParse(consumer)
}

func (c TCPClient) readAndParse(consumer chan<- file.FileContent) {
	for { // todo : there should be a way to stop the consumer
		bufferFileName := make([]byte, utils.CommandLength)
		bufferFileSize := make([]byte, utils.CommandLength)

		_, fNameErr := c.socket.Read(bufferFileName)
		c.handleError(fNameErr)
		_, fSizeError := c.socket.Read(bufferFileSize)
		c.handleError(fSizeError)

		fileName := strings.Trim(string(bufferFileName), utils.Filler)
		fileSize, convErr := strconv.ParseInt(strings.Trim(string(bufferFileSize), utils.Filler), 10, 64)
		c.handleError(convErr)

		fileBuffer := make([]byte, fileSize)
		_, fileErr := c.socket.Read(fileBuffer) //todo : Would have to read as small buffers and keep appending
		c.handleError(fileErr)

		fileContent := file.NewFileContent(fileSize, fileName)
		fileContent.Append(fileBuffer)
		consumer <- fileContent
		//todo : status has to be reported and persisted. This should be used as latest status while reconnecting
	}
}

func (c TCPClient) handleError(err error) {
	if err != nil {
		fmt.Println("errored", err)
		//todo : ERROR HAS TO BE HANDLED BY FETCHING THE LATEST STATUS AND RECONNECTING
	}
}

func (c TCPClient) Disconnect() {
	c.socket.Close()
}
