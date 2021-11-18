package server

import "net/http"

type Server interface {
	Run()
}

type edgeServer struct {
	hubServer *http.Server
	proxyServer *http.Server
	secureProxyServer *http.Server

	stopChan <-chan struct{}
}

func CreateServer() *edgeServer {
	return &edgeServer{
		hubServer: nil,
		proxyServer: nil,
		secureProxyServer: nil,
	}
}

