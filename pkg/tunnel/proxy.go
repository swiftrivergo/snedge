package tunnel

import (
	"github.com/swiftrivergo/snedge/pkg/util"
	"net/http"
)

type Proxy struct {
	server http.Server
	Tunnel Tunnel
}

func NewProxy() *Proxy {
	return new(Proxy)
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	util.HandleHTTP(rw, req)
}

func (p *Proxy) ServeTunnel(w http.ResponseWriter, req *http.Request) {
	util.HandleTunnel(w, req)
}

func (p *Proxy) Listen() error {
	p.Tunnel = New()
	return p.Tunnel.Listen()
}
