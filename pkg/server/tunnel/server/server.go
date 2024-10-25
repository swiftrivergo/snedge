package server

import (
	"crypto/subtle"
	"fmt"
	"github.com/rancher/remotedialer"
	"github.com/swiftrivergo/snedge/pkg/server/tunnel/router"
	"github.com/swiftrivergo/snedge/pkg/server/tunnel/transport"
	"github.com/twinj/uuid"
	"k8s.io/apimachinery/pkg/util/proxy"
	"log"
	"net/http"
	"sync"
)

// Server for the exit-server of ynCmd
type Server struct {

	// Port serves data to clients
	Port int

	// ControlPort represents the tunnel to the ynCmd client
	ControlPort int

	BindAddr string

	// Token is used to authenticate a client
	Token string

	router router.Router
	server *remotedialer.Server

	// DisableWrapTransport prevents CORS headers from being striped from responses
	DisableWrapTransport bool
}

// Serve traffic
func (s *Server) Serve() {
	log.Printf("BindAddr: [%s] Port: [%d] ControlPort:[%d]", s.BindAddr, s.Port, s.ControlPort)

	if s.ControlPort == s.Port {
		s.server = remotedialer.New(s.authorized, remotedialer.DefaultErrorWriter)
		s.router.Server = s.server

		//http.HandleFunc("/", util.HandleHTTP)
		//http.HandleFunc("/tunnel", util.HandleTunnel)

		server := http.NewServeMux()
		server.HandleFunc("/", s.proxy)
		server.HandleFunc("/tunnel", s.tunnel)

		log.Printf("(1) Data Plane Listening on %s:%d:%d\n", s.BindAddr, s.ControlPort, s.Port)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.BindAddr, s.Port), server); err != nil {
			log.Fatal(err)
		}
	} else {

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()

			controlServer := http.NewServeMux()
			s.server = remotedialer.New(s.authorized, remotedialer.DefaultErrorWriter)
			s.router.Server = s.server

			controlServer.HandleFunc("/", s.proxy)
			//controlServer.HandleFunc("/tunnel", s.tunnel)

			log.Printf("(2) Data Plane Listening on %s:%d\n", s.BindAddr, s.Port)
			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.BindAddr, s.Port), controlServer); err != nil {
				log.Fatal(err)
			}

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			controlServer := http.NewServeMux()
			controlServer.HandleFunc("/", s.proxy)

			http.HandleFunc("/", s.proxy)
			log.Printf("(2) Control Plane Listening on %s:%d\n", s.BindAddr, s.ControlPort)

			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.BindAddr, s.ControlPort), controlServer); err != nil {
				log.Fatal(err)
			}
		}()

		wg.Wait()
	}
}

func (s *Server) tunnel(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
	s.router.Remove(r)
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	route := s.router.Lookup(r)
	if route == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	ynCmdID := uuid.Formatter(uuid.NewV4(), uuid.FormatHex)
	log.Printf("[%s] proxy %s %s %s", ynCmdID, r.Host, r.Method, r.URL.String())
	r.Header.Set(transport.InHeader, ynCmdID)

	u := *r.URL
	u.Host = r.Host
	u.Scheme = route.Scheme

	httpProxy := proxy.NewUpgradeAwareHandler(&u, route.Transport, !s.DisableWrapTransport, false, s)
	httpProxy.ServeHTTP(w, r)
}

func (s *Server) Error(w http.ResponseWriter, req *http.Request, err error) {
	remotedialer.DefaultErrorWriter(w, req, http.StatusInternalServerError, err)
}

func (s *Server) tokenValid(req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	return len(s.Token) == 0 || subtle.ConstantTimeCompare([]byte(auth), []byte("Bearer "+s.Token)) == 1
}

func (s *Server) authorized(req *http.Request) (id string, ok bool, err error) {
	defer func() {
		if id == "" {
			// empty id is also an auth failure
			ok = false
		}
		if !ok || err != nil {
			// don't let non-authed request clear routes
			req.Header.Del(transport.InHeader)
		}
	}()

	if !s.tokenValid(req) {
		return "", false, nil
	}

	return s.router.Add(req), true, nil
}

