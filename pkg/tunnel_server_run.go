package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/swiftrivergo/snedge/pkg/proxy"
	"github.com/swiftrivergo/snedge/pkg/tunnel"
	"log"
	"net/http"
)

func main() {
	var pemPath string
	flag.StringVar(&pemPath, "pem", "d:\\ssh\\server.pem", "path to pem file")
	var keyPath string
	flag.StringVar(&keyPath, "key", "d:\\ssh\\server.key", "path to key file")
	var proto string
	flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
	flag.Parse()

	if proto != "http" && proto != "https" {
		//log.Fatal("Protocol must be either http or https")
		fmt.Println("Protocol must be either http or https")
	}
	server := &http.Server{
		Addr: ":8082",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				tunnel.HandleTunnel(w, r)
			} else {
				tunnel.HandleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	//v2:
	var tl proxy.Tunnel
	tu := proxy.New()
	tu.SetAddr(server.Addr)
	fmt.Println(server.Addr, tu.GetAddr())
	tl = tu

	if proto == "http" {
		go func(t proxy.Tunnel) {
			//Todo: addr should be set by user
			log.Fatal(tl.Listen())
		}(tl)
	} else {
		go func(t proxy.Tunnel) {
			//Todo: TLS should be support; addr should be set by user
			log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
		}(tl)
	}

	select {
	}
}
