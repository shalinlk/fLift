package client

import (
	"github.com/shalinlk/fLift/file"
)

type Client interface {
	Register()
	Start(chan<- file.FileContent)
	Disconnect()
}
