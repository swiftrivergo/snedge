package tunnel

import (
	"crypto/tls"
	"github.com/swiftrivergo/snedge/pkg/server"
	"net/http"
)

type Middleware interface {
	WrapHandler(http.Handler) http.Handler
	Name() string
}
type HandlerWrappers []Middleware

type tunnelServer struct {
	//egressSelectorEnable     bool
	//interceptorServerUDSFile string
	masterAddr               string
	masterInsecureAddr       string
	agentAddr                string
	serverCount              int
	tlsCfg                   *tls.Config
	wrappers                 []Middleware
	//proxyStrategy            string
}

func newTunnelServer() server.TunnelServer {
	var ts server.TunnelServer = &tunnelServer{}
	return ts
}

func (t *tunnelServer) Run() error {
	return nil
}