package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shalinlk/fLift/file"
	"github.com/shalinlk/fLift/utils"
)

type TCPClient struct {
	clientId       string
	serverHost     string
	consumerChan   chan file.FileContent
	connectionPool []net.Conn
	poolLock       *sync.Mutex
}

func NewTCPClient(serverHost string, consumerChan chan file.FileContent, concurrentConnections int) TCPClient {
	client := TCPClient{
		serverHost:     serverHost,
		consumerChan:   consumerChan,
		connectionPool: make([]net.Conn, concurrentConnections),
		poolLock:       &sync.Mutex{},
	}
	connectionCount := 0
	for i := 0; i < concurrentConnections; i++ {
		conn := client.dial()
		if conn != nil {
			client.connectionPool[connectionCount] = conn
			connectionCount++
		}
	}
	return client
}

func (c *TCPClient) dial() net.Conn {
	fmt.Println("Connecting to ", c.serverHost)
	socket, err := net.Dial("tcp", c.serverHost)
	if err != nil {
		fmt.Println("socket failed to " + c.serverHost)
		return nil
	}
	return socket
}

func (c *TCPClient) redial(connectionIndex int) {
	socket, err := net.Dial("tcp", c.serverHost)
	if err != nil {
		c.redial(0)
	} else {
		c.poolLock.Lock()
		//todo : add to specific location
		c.poolLock.Unlock()
		c.readAndParse(socket, connectionIndex)
	}
}

//func (c TCPClient) Register() {
//	_, _ = c.socket.Write([]byte(utils.FillUpForCommand("REGISTER")))
//	clientIdBuffer := make([]byte, 60)
//	_, _ = c.socket.Read(clientIdBuffer)
//	c.clientId = strings.Trim(string(clientIdBuffer), ":")
//}

func (c TCPClient) Start() {
	//_, _ = c.socket.Write([]byte(utils.FillUpForCommand("START")))
	if len(c.connectionPool) == 0 {
		panic("could not connect any client to producer")
	}
	for i := 0; i < len(c.connectionPool); i++ {
		go c.readAndParse(c.connectionPool[i], i)
	}
	time.Sleep(time.Second * 10)
}

func (c TCPClient) readAndParse(conn net.Conn, index int) {
	for { // todo : there should be a way to stop the consumer
		bufferFileName := make([]byte, utils.CommandLength)
		bufferFileSize := make([]byte, utils.CommandLength)
		bufferFilePath := make([]byte, utils.CommandLength)

		_, fNameErr := conn.Read(bufferFileName)
		c.handleError(fNameErr, index)

		_, fSizeError := conn.Read(bufferFileSize)
		c.handleError(fSizeError, index)

		_, fPathErr := conn.Read(bufferFilePath)
		c.handleError(fPathErr, index)

		fileName := strings.Trim(string(bufferFileName), utils.Filler)
		fileSize, convErr := strconv.Atoi(strings.Trim(string(bufferFileSize), utils.Filler))
		c.handleError(convErr, index)
		filePath := strings.Trim(string(bufferFilePath), utils.Filler)

		fileBuffer := make([]byte, fileSize)
		_, fileErr := conn.Read(fileBuffer) //todo : Would have to read as small buffers and keep appending
		c.handleError(fileErr, index)

		fileContent := file.NewFileContent(fileSize, fileName, 0, filePath, fileBuffer)
		c.consumerChan <- fileContent
		//todo : status has to be reported and persisted. This should be used as latest status while reconnecting
	}
}

func (c TCPClient) handleError(err error, connectionIndex int) {
	if err != nil {
		fmt.Println("errored", err)
		//todo : ERROR HAS TO BE HANDLED BY FETCHING THE LATEST STATUS
		c.poolLock.Lock()
		//todo : remove connection from pool
		c.poolLock.Unlock()
		c.redial(connectionIndex)
	}
}
