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

	p := tunnel.NewProxy()

	//v2:
	var tl tunnel.Tunnel
	tu := tunnel.New()
	tu.SetAddr(server.Addr)
	log.Println("server start:", server.Addr, "listen GetListenAddr():", tu.GetListenAddr())
	tl = tu
	p.Tunnel = tl

	if proto == "http" {
		go func(t tunnel.Tunnel) {
			//Todo: addr should be set by user
			log.Fatal(p.Tunnel.Listen())
		}(tl)
	} else {
		go func(t tunnel.Tunnel) {
			//Todo: TLS should be support; addr should be set by user
			log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
		}(tl)
	}

	select {
	}
}
