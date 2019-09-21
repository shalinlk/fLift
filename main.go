package main

import (
	"flag"
	"fmt"
	. "github.com/shalinlk/fLift/file"
	. "github.com/shalinlk/fLift/server"
	"runtime"
	"time"

	"github.com/shalinlk/fLift/client"
	"github.com/shalinlk/fLift/utils"
)

const ModeConsumer = "consumer"
const ModeProducer = "producer"
const OperationModeStart = "start"
const OperationModeRestart = "restart"

func main() {
	env := flag.String("env", "", "environment")
	mode := flag.String("mode", ModeConsumer, ModeConsumer+" / "+ModeProducer)
	operationMode := flag.String("operationMode", OperationModeRestart, OperationModeStart + " / " + OperationModeRestart)
	flag.Parse()

	if *mode != ModeConsumer && *mode != ModeProducer {
		panic("mode should be " + ModeConsumer + " / " + ModeProducer)
	}

	fmt.Println("Welcome to fLift; running in " + *mode + " mode")
	if *mode  == ModeProducer{
		fmt.Println("Operation Mode : ", *operationMode)
	}
	fmt.Println("Starting time : ", time.Now().UnixNano() / int64(time.Millisecond))

	config := utils.ReadConfig(*env)

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *mode == ModeConsumer {
		client.StartConsumer(
			config.Host,
			config.WriteFilePath,
			config.ConnectionType,
			config.WriteBufferSize,
			config.WriterCount)
	} else if *mode == ModeProducer {
		reader := NewReader(config.ReadFilePath, config.ReadBufferSize, config.ReaderCount, config.ReadBatchSize)
		server := NewServer(
			config.Port,
			reader,
			*operationMode,
			config.MaxClients,
			config.StatusFlushInterval,
			config.KeepStatus)
		server.Start()
	}
}
