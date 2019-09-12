package client

type Client interface {
	Register()
	Start()
	Disconnect()
}
