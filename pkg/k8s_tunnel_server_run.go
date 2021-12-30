package main

import (
	svr "github.com/swiftrivergo/snedge/pkg/server/tunnel/server"
)

func main() {

	s := &svr.Server{
		Port:                 8082,
		ControlPort:          8083,
		BindAddr:             "",
		Token:                "",
		DisableWrapTransport: false,
	}

	s.Serve()
}
