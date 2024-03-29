package server

import (
	"fmt"
	. "github.com/shalinlk/fLift/file"
	"github.com/shalinlk/fLift/utils"
	"net"
	"strconv"
	"time"
)

type Server struct {
	connections        map[string]net.Conn
	contentProducer    chan FileContent
	connectionProducer chan net.Conn
	port               int
	reader             Reader
	operationMode      string
	statusTracker      StatusTracker
}

func NewServer(port int, reader Reader, operationMode string, maxClients, statusFlushInterval int, keepStatus bool) Server {
	connections := make(map[string]net.Conn)
	connectionProducer := make(chan net.Conn)
	statusTracker := NewStatusTracker(maxClients, statusFlushInterval, operationMode, keepStatus)
	return Server{
		connections:        connections,
		contentProducer:    reader.Feeder,
		connectionProducer: connectionProducer,
		port:               port,
		reader:             reader,
		operationMode:      operationMode,
		statusTracker:      statusTracker,
	}
}

func (s Server) Start() {
	go s.reader.Start(s.operationMode, s.statusTracker.CurrentIndex())
	s.statusTracker.Start()
	s.acceptConnection(s.port)
}

func (s Server) acceptConnection(port int) {
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	panicOnErrorWithMessage(err, "could not start server")
	defer server.Close()
	for {
		connection, connError := server.Accept()
		if connError != nil {
			fmt.Println("error in accepting connection", connError)
		} else {
			fmt.Println("\nAccepted connection at: ", time.Now().UnixNano() / int64(time.Millisecond))
			go s.feeder(connection)
		}
	}
}

func (s Server) addConnectionToPool(conn net.Conn) {
	s.connections[fmt.Sprintf("%d", len(s.connections)+1)] = conn
}

func (s Server) feeder(conn net.Conn) {
	trackerChan := s.statusTracker.StatusTrackerChan()
	for content := range s.contentProducer{
		fullLengthName, err := utils.FillUpForCommand(content.Name, utils.NameLength)
		if err != nil {
			fmt.Println("Failed to format " + content.Name + "; Size exceeds; Skipping file :", content.Name)
			continue
		}
		fullLengthSize, err := utils.FillUpForCommand(strconv.Itoa(content.Size), utils.SizeLength)
		if err != nil {
			fmt.Println("Failed to format " + strconv.Itoa(content.Size) + "; Size exceeds; Skipping file :", content.Name)
			continue
		}
		fullLengthPath, err := utils.FillUpForCommand(content.Path, utils.PathLength)
		if err != nil {
			fmt.Println("Failed to format " + content.Path + "; Size exceeds; Skipping file :", content.Name)
			continue
		}

		_, err = conn.Write([]byte(fullLengthName))
		if err != nil {
			fmt.Println("Writing file name to connection failed. Dropping connection")
			return
		}

		_, err = conn.Write([]byte(fullLengthSize))
		if err != nil {
			fmt.Println("Writing file size to connection failed. Dropping connection")
			return
		}

		_, err = conn.Write([]byte(fullLengthPath))
		if err != nil {
			fmt.Println("Writing file path to connection failed. Dropping connection")
			return
		}

		_, err = conn.Write(content.Content)
		if err != nil {
			fmt.Println("Writing file content to connection failed. Dropping connection")
			return
		}
		trackerChan <- content.Index
	}

	fmt.Println("writing eof on connection at : ", time.Now().UnixNano() / int64(time.Millisecond))
	tempFileName, err := utils.FillUpForCommand("eof.temp", utils.NameLength)
	if err != nil {
		fmt.Println("Failed to format eof.temp; Size exceeds; Skipping file :")
		return
	}
	tempFileSize, err := utils.FillUpForCommand(strconv.Itoa(-1), utils.SizeLength)
	if err != nil {
		fmt.Println("Failed to format -1; Size exceeds; Skipping file eof.temp:")
		return
	}
	_, err = conn.Write([]byte(tempFileName))
	_, err = conn.Write([]byte(tempFileSize))
}

func panicOnErrorWithMessage(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
