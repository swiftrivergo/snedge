package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

type Tunnel interface {
	Listen() error
}

type tunnel struct {
	TCP string
	Port string
}

func New() *tunnel {
	tl := new(tunnel)
	tl.TCP = "tcp"
	tl.Port = "80"
	return tl
}

func (t *tunnel) Listen() error {
	l, err := net.Listen(t.TCP, ":" + t.Port)

	if err != nil {
		log.Panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		fmt.Println()
		go handleClient(c)
	}
	return nil
}

func handleClient(c net.Conn) {
	if c == nil {
		return
	}
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	log.Println("local:", c.LocalAddr(),"remote:", c.RemoteAddr())

	var b [1024]byte
	n, err := c.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	sscanf, err := fmt.Sscanf(string(b[:bytes.IndexByte(b[:],'\n')]), "%s%s", &method, &host)
	if err != nil {
		log.Println(sscanf, err)
		return
	}

	fmt.Println("method:", method, "host:", host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}
	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			//Todo port should be passed by user
			address = hostPortURL.Host + ":8081"
		} else {
			address = hostPortURL.Host
		}
	}

	fmt.Println("address:", address)
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" || method == "connect" {
		fprint, err := fmt.Fprint(c, "HTTP/1.1 200 Connection establish\r\n")
		if err != nil {
			log.Println(fprint, err)
			return
		}
	} else {
		write, err := server.Write(b[:n])
		if err != nil {
			log.Println(write, err)
			return
		}
	}

	go func() {
		_, err := io.Copy(server, c)
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = io.Copy(c, server)
	if err != nil {
		log.Println(err)
		return
	}
}
