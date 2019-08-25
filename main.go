package main

import (
	"flag"
	"fmt"

	"github.com/shalinlk/fLift/client"
	"github.com/shalinlk/fLift/utils"
)

func main() {
	fmt.Println("welcome to fLift")
	env := flag.String("env", "", "environment")
	mode := flag.String("mode", "consumer", "conuser / producer")

	flag.Parse()

	conifg := utils.ReadConfig(*env)

	if *mode == "consuemer" {
		client.StartConsumer(conifg.Host, conifg.BaseFilePath, conifg.ConnectionType, conifg.WriteBufferSize)
	}
}
