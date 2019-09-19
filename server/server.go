package server

import (
	"fmt"
	. "github.com/shalinlk/fLift/file"
	"github.com/shalinlk/fLift/utils"
	"net"
	"strconv"
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
			go s.feeder(connection)
		}
	}
}

func (s Server) addConnectionToPool(conn net.Conn) {
	s.connections[fmt.Sprintf("%d", len(s.connections)+1)] = conn
}

func (s Server) feeder(conn net.Conn) {
	var err error
	trackerChan := s.statusTracker.StatusTrackerChan()
	for {
		content := <-s.contentProducer

		_, err = conn.Write([]byte(utils.FillUpForCommand(content.Name, utils.NameLength)))
		if err != nil {
			fmt.Println("Writing file name to connection failed. Dropping connection")
			return
		}

		_, err = conn.Write([]byte(utils.FillUpForCommand(strconv.Itoa(content.Size), utils.SizeLength)))
		if err != nil {
			fmt.Println("Writing file size to connection failed. Dropping connection")
			return
		}

		_, err = conn.Write([]byte(utils.FillUpForCommand(content.Path, utils.PathLength)))
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
}

func panicOnErrorWithMessage(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
