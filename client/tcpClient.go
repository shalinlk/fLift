package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shalinlk/fLift/file"
	"github.com/shalinlk/fLift/utils"
)

type TCPClient struct {
	clientId     string
	socket       net.Conn
	serverHost   string
	consumerChan chan file.FileContent
}

func NewTCPClient(serverHost string, consumerChan chan file.FileContent) *TCPClient {
	client := TCPClient{
		serverHost:   serverHost,
		consumerChan: consumerChan,
	}
	client.dial()
	return &client
}

func (c *TCPClient) dial() {
	fmt.Println("Connecting to ", c.serverHost)
	socket, err := net.Dial("tcp", c.serverHost)
	if err != nil {
		fmt.Println("socket failed to " + c.serverHost)
		panic(err)
	}
	c.socket = socket
}

func (c *TCPClient) redial() {
	socket, err := net.Dial("tcp", c.serverHost)
	if err != nil {
		c.redial()
	}else {
		c.socket = socket
		c.readAndParse()
	}
}

func (c *TCPClient) Register() {
	filledUpRegCommand, _ := utils.FillUpForCommand("REGISTER", utils.SizeLength)
	_, _ = c.socket.Write([]byte(filledUpRegCommand))
	clientIdBuffer := make([]byte, 60)
	_, _ = c.socket.Read(clientIdBuffer)
	c.clientId = strings.Trim(string(clientIdBuffer), ":")
}

func (c *TCPClient) Start() {
	//_, _ = c.socket.Write([]byte(utils.FillUpForCommand("START")))
	c.readAndParse()
	close(c.consumerChan)
	c.Disconnect()
}

func (c *TCPClient) readAndParse() {
	for { // todo : there should be a way to stop the consumer
		bufferFileName := make([]byte, utils.NameLength)
		bufferFileSize := make([]byte, utils.SizeLength)
		bufferFilePath := make([]byte, utils.PathLength)

		_, fNameErr := c.socket.Read(bufferFileName)
		c.reconnectOnError(fNameErr)

		_, fSizeError := c.socket.Read(bufferFileSize)
		c.reconnectOnError(fSizeError)

		fileName := strings.Trim(string(bufferFileName), utils.Filler)
		fileSize, convErr := strconv.Atoi(strings.Trim(string(bufferFileSize), utils.Filler))
		if fileSize < 0 {
			fmt.Println("EOF file read; Stopping to read from conn at : ", time.Now().UnixNano() / int64(time.Millisecond))
			return
		}

		_, fPathErr := c.socket.Read(bufferFilePath)
		c.reconnectOnError(fPathErr)

		c.reconnectOnError(convErr)
		filePath := strings.Trim(string(bufferFilePath), utils.Filler)

		fileBuffer := make([]byte, fileSize)
		_, fileErr := c.socket.Read(fileBuffer) //todo : Would have to read as small buffers and keep appending
		c.reconnectOnError(fileErr)

		fileContent := file.NewFileContent(fileSize, fileName, 0, filePath, fileBuffer)
		c.consumerChan <- fileContent
		//todo : status has to be reported and persisted. This should be used as latest status while reconnecting
	}
}

func (c *TCPClient) reconnectOnError(err error) {
	if err != nil {
		fmt.Println("errored on connection : ", err)
		//todo : ERROR HAS TO BE HANDLED BY FETCHING THE LATEST STATUS
		c.redial()
	}
}

func (c *TCPClient) Disconnect() {
	_ = c.socket.Close()
}
