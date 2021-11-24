package ynhub

import (
	"crypto/tls"
	"github.com/gorilla/mux"
	"github.com/swiftrivergo/snedge/pkg/app/config"
	"github.com/swiftrivergo/snedge/pkg/server"
	"net"
	"net/http"
)

type ProxyServer interface {
	ListenAndServe() error
}

type SecureProxyServer interface {
	ListenAndServeTLS(certFile, keyFile string) error
}

type ynEdgeHubServer struct {
	hub              ProxyServer
	proxy            ProxyServer
	secureProxy      SecureProxyServer
	dummyProxy       ProxyServer
	dummySecureProxy SecureProxyServer
}

type hubServer struct {
	Addr           string
	ProxyHandler   http.Handler
	MaxHeaderBytes int //default 1 << 20
	*http.Server
}

type proxyServer struct {
	Addr         string
	ProxyHandler http.Handler
	*http.Server
}

type secureProxyServer struct {
	SecureAddr     string
	ProxyHandler   http.Handler
	TTLSConfig     *tls.Config
	TLSNextProto   map[string]func(*http.Server, *tls.Conn, http.Handler)
	MaxHeaderBytes int //default 1 << 20
	*http.Server
}

func newHubServer(addr string, handler http.Handler) ProxyServer {
	hub := &hubServer{}
	hub.Addr = addr
	hub.Handler = handler
	hub.MaxHeaderBytes = 1 << 20

	return hub
}

func newProxyServer(addr string, handler http.Handler) ProxyServer {
	proxy := &proxyServer{}
	proxy.Addr = addr
	proxy.Handler = handler

	return proxy
}

func newSecureProxyServer(addr string, handler http.Handler, config *tls.Config) SecureProxyServer {
	secureProxy := &secureProxyServer{}
	secureProxy.Addr = addr
	secureProxy.Handler = handler
	secureProxy.TLSConfig = config
	secureProxy.MaxHeaderBytes = 1 << 20
	secureProxy.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))

	return secureProxy
}

func newDummyProxyServer(addr string, handler http.Handler) ProxyServer {
	proxy := &proxyServer{}
	proxy.Addr = addr
	proxy.Handler = handler
	proxy.MaxHeaderBytes = 1 << 20
	return proxy
}

func newDummySecureProxyServer(addr string, handler http.Handler, config *tls.Config) SecureProxyServer {
	secureProxy := &secureProxyServer{}
	secureProxy.Addr = addr
	secureProxy.Handler = handler
	secureProxy.MaxHeaderBytes = 1 << 20
	secureProxy.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	return secureProxy
}

func (h *hubServer) ListenAndServe() error {
	return h.ListenAndServe()
}

func (p *proxyServer) ListenAndServe() error {
	return p.ListenAndServe()
}

func (s *secureProxyServer) ListenAndServeTLS(certFile, keyFile string) error {
	return s.ListenAndServeTLS(certFile, keyFile)
}

func NewYnEdgeHubServer(cfg *config.EdgeHubConfig, proxyHandel http.Handler) (server.HubServer, error) {
	hubMux := mux.NewRouter()
	registerHandlers(hubMux, cfg)

	var hubSvr ProxyServer
	var proxySvr ProxyServer
	var dummyProxySvr ProxyServer
	var dummySecureProxySvr SecureProxyServer

	hubSvr = newHubServer(cfg.HubServerAddr, hubMux)
	proxySvr = newProxyServer(cfg.ProxyServerAddr, proxyHandel)

	secureProxySvr := newSecureProxyServer(cfg.SecureProxyServerAddr, proxyHandel, cfg.TLSConfig)
	if cfg.EnableDummyIf {
		if _, err := net.InterfaceByName(cfg.HubAgentDummyIfName); err != nil {
			return nil, err
		}

		dummyProxySvr = newDummyProxyServer(cfg.DummyProxyServerAddr, proxyHandel)
		dummySecureProxySvr = newDummySecureProxyServer(cfg.DummySecureProxyServerAddr, proxyHandel, cfg.TLSConfig)
	}

	return &ynEdgeHubServer{
		hub:              hubSvr,
		proxy:            proxySvr,
		secureProxy:      secureProxySvr,
		dummyProxy:       dummyProxySvr,
		dummySecureProxy: dummySecureProxySvr,
	}, nil
}

func (s *ynEdgeHubServer) Run() {
	go func() {
		err := s.hub.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	if s.dummyProxy != nil {
		go func() {
			err := s.dummyProxy.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
		go func() {
			err := s.dummySecureProxy.ListenAndServeTLS("", "")
			if err != nil {
				panic(err)
			}
		}()
	}

	go func() {
		err := s.secureProxy.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	}()

	err := s.proxy.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func registerHandlers(hubMux *mux.Router, cfg *config.EdgeHubConfig) {
	//Todo

}
