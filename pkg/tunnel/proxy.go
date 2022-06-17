package tunnel

import (
	"github.com/swiftrivergo/snedge/pkg/util"
	"net/http"
)

type Proxy struct {
	server *http.Server
	Tunnel Tunnel
}

func NewProxy() *Proxy {
	t := New()
	s := new(http.Server)

	proxy := new(Proxy)
	proxy.server = s
	proxy.Tunnel = t

	return proxy
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	util.HandleHTTP(rw, req)
}

func (p *Proxy) ServeTunnel(w http.ResponseWriter, req *http.Request) {
	util.HandleTunnel(w, req)
}

func (p *Proxy) Run() error {
	return p.Tunnel.Run()
}

func (p *Proxy) GetServer() *http.Server {
	return p.server
}

func (p *Proxy) SetServer(s *http.Server) {
	p.server = s
	p.setAddr(p.server.Addr)
}

func (p *Proxy) setAddr(addr string) {
	switch p.Tunnel.(type) {
	case *tunnel:
		tl := p.Tunnel.(*tunnel)
		//Todo: net.Listener should be support.
		tl.AddListenAddr(addr)
	default:
	}
}

func (p *Proxy) SetTunnel(t Tunnel) {
	p.Tunnel = t
	p.setAddr(p.server.Addr)
}
