package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Host            string
	ConnectionType  string
	WriteBufferSize int
	BaseFilePath    string
}

func ReadConfig(env string) Config {
	if env != "" {
		env = env + "-"
	}
	fileContent, err := ioutil.ReadFile("/" + env + "properties.json")
	if err != nil {
		fmt.Println("Could not load config file ")
		panic(err)
	}
	config := Config{}
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		fmt.Println("Could not parse properties file")
	}
	return config
}
