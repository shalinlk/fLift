package server

import (
	"fmt"
	"github.com/shalinlk/fLift/file"
	"github.com/shalinlk/fLift/utils"
	"net"
	"strconv"
)

type Server struct {
	connections        map[string]net.Conn
	contentProducer    chan file.FileContent
	connectionProducer chan net.Conn
	port               int
	reader             file.Reader
}

func NewServer(port int, reader file.Reader) Server {
	connections := make(map[string]net.Conn)
	connectionProducer := make(chan net.Conn)
	return Server{
		connections:        connections,
		contentProducer:    reader.Feeder,
		connectionProducer: connectionProducer,
		port:               port,
		reader:             reader,
	}
}

func (s Server) Start() {
	//go s.feeder()
	go s.reader.Start()
	s.acceptConnection(s.port)
}

func (s Server) acceptConnection(port int) {
	server, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	panicOnErrorWithMessage(err, "could not start server")
	defer server.Close()
	for {
		connection, connError := server.Accept()
		panicOnErrorWithMessage(connError, "error in accepting connection")
		go s.feeder(connection)
	}
}


func (s Server) addConnectionToPool(conn net.Conn) {
	s.connections[fmt.Sprintf("%d", len(s.connections)+1)] = conn
}

func (s Server) feeder(conn net.Conn) {
	//var err error
	for {
		select {
		case content := <- s.contentProducer:
			_, _ = conn.Write([]byte(utils.FillUpForCommand(content.Name)))
			_, _ = conn.Write([]byte(utils.FillUpForCommand(strconv.Itoa(content.Size))))
			_, _ = conn.Write(content.Content)
		}
	}
}

func panicOnErrorWithMessage(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
