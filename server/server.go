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
	go s.reader.Start()
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
		}else {
			go s.feeder(connection)
		}
	}
}

func (s Server) addConnectionToPool(conn net.Conn) {
	s.connections[fmt.Sprintf("%d", len(s.connections)+1)] = conn
}

func (s Server) feeder(conn net.Conn) {
	var err error
	for {
		content := <-s.contentProducer
		_, err = conn.Write([]byte(utils.FillUpForCommand(content.Name)))
		if err != nil {
			fmt.Println("Writing file name to connection failed. Dropping connection")
			return
		}
		_, err = conn.Write([]byte(utils.FillUpForCommand(strconv.Itoa(content.Size))))
		if err != nil {
			fmt.Println("Writing file size to connection failed. Dropping connection")
			return
		}
		_, err = conn.Write(content.Content)
		if err != nil {
			fmt.Println("Writing file content to connection failed. Dropping connection")
			return
		}
	}
}

func panicOnErrorWithMessage(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
