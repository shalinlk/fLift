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
	go s.feeder()
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
		s.connectionProducer <- connection
	}
}

func (s Server) feeder() {
	for {
		select {
		case conn := <-s.connectionProducer:
			s.addConnectionToPool(conn)
		case fileContent := <-s.contentProducer:
			s.sendContent(fileContent)
		}
	}
}

func (s Server) addConnectionToPool(conn net.Conn) {
	s.connections[fmt.Sprintf("%d", len(s.connections)+1)] = conn
}

func (s Server) sendContent(file file.FileContent) {
	for _, v := range s.connections {
		_, _ = v.Write([]byte(utils.FillUpForCommand(file.Name)))
		_, _ = v.Write([]byte(utils.FillUpForCommand(strconv.Itoa(file.Size))))
		_, _ = v.Write(file.Content)
	}
}

func panicOnErrorWithMessage(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
