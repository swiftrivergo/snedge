package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/swiftrivergo/snedge/pkg/tunnel"
	"github.com/swiftrivergo/snedge/pkg/util"
	"log"
	"net/http"
)

func main() {
	var pemPath string
	flag.StringVar(&pemPath, "pem", "d:\\ssh\\server.pem", "path to pem file")
	var keyPath string
	flag.StringVar(&keyPath, "key", "d:\\ssh\\server.key", "path to key file")
	var proto string
	flag.StringVar(&proto, "proto", "http", "Proxy protocol (http or https)")
	flag.Parse()

	if proto != "http" && proto != "https" {
		//log.Fatal("Protocol must be either http or https")
		fmt.Println("Protocol must be either http or https")
	}
	server := &http.Server{
		Addr: ":8082",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				util.HandleTunnel(w, r)
			} else {
				util.HandleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	//v2:
	p := tunnel.NewProxy()
	tu := tunnel.New()
	tu.SetForwardPort("8081")
	tu.Addr = server.Addr
	p.Tunnel = tu

	p.SetServer(server)
	p.SetAddr(server.Addr)

	log.Println("server start:", tu.Addr, "listen GetListenAddr():", tu.GetListenAddr())

	if proto == "http" {
		go func(t tunnel.Tunnel) {
			//Todo: addr should be set by user
			log.Fatal(t.Listen())
		}(p.Tunnel)
	} else {
		go func(t tunnel.Tunnel) {
			//Todo: TLS should be supportted; addr should be set by user
			log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
		}(tu)
	}

	select {
	}
}
