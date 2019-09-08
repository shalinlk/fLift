package main

import (
	"flag"
	"fmt"
	. "github.com/shalinlk/fLift/file"
	. "github.com/shalinlk/fLift/server"
	"runtime"

	"github.com/shalinlk/fLift/client"
	"github.com/shalinlk/fLift/utils"
)

const ModeConsumer = "consumer"
const ModeProducer = "producer"
	
func main() {
	env := flag.String("env", "", "environment")
	mode := flag.String("mode", ModeConsumer, ModeConsumer + " / " + ModeProducer)
	flag.Parse()

	if *mode != ModeConsumer && *mode != ModeProducer {
		panic("mode should be " + ModeConsumer + " / " + ModeProducer)
	}

	fmt.Println("welcome to fLift; running in " + *mode + " mode")

	config := utils.ReadConfig(*env)

	if *mode == ModeConsumer {
		runtime.GOMAXPROCS(runtime.NumCPU())
		client.StartConsumer(config.Host, config.WriteFilePath, config.ConnectionType, config.WriteBufferSize, config.WriterCount)
	} else if *mode == ModeProducer {
		reader := NewReader(config.ReadFilePath, config.ReadBufferSize)
		server := NewServer(config.Port, reader)
		server.Start()
	}
}
