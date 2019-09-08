package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	ConnectionType  string `json:"connection_type"`
	WriteBufferSize int    `json:"write_buffer_size"`
	WriteFilePath   string `json:"write_file_path"`
	WriterCount     int    `json:"writer_count"`
	ReadBufferSize  int    `json:"read_buffer_size"`
	ReadFilePath    string `json:"read_file_path"`
}

func ReadConfig(env string) Config {
	if env != "" {
		env = env + "-"
	}
	fileContent, err := ioutil.ReadFile(env + "properties.json")
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
