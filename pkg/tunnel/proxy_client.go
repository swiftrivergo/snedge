package tunnel

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

type ListenFunc func(string) error

type tunnel struct {
	tcp  string
	Addr string

	dialAddr string
	dialPort string


	ListenHost string
	ListenPort string

	method string
	host string
	hostName string
	hostPort string

	local string
	remote string
}

func New() *tunnel {
	tl := new(tunnel)
	tl.tcp = "tcp"
	return tl
}

func (t *tunnel) SetListenPort(port string) {
	t.ListenPort = port
}

func (t *tunnel) GetListenAddr() string {
	return t.ListenHost + ":" + t.ListenPort
}

func (t *tunnel) SetAddr(address string) {
	t.dialAddr = strings.Trim(address, " ")
	if strings.Index(t.dialAddr, ":") == -1 {
		t.ListenHost = ":80"
	} else {
		t.ListenHost = strings.Split(t.dialAddr, ":")[0]
		t.ListenPort = strings.SplitAfter(strings.Trim(t.dialAddr, " "), ":")[1]
		log.Println(address, t.dialAddr, t.tcp, t.ListenHost, t.ListenPort)
	}
}

func (t *tunnel) Listen() error {
	l, err := net.Listen(t.tcp, t.GetListenAddr())

	log.Println("tunnel listen tcp:", t.tcp, "GetListenAddr():", t.GetListenAddr(), "ListenPort:", t.ListenPort)
	if err != nil {
		log.Panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClient(c, "8081")
	}
	return nil
}

func handleClient(c net.Conn, forwardPort string) {
	if c == nil {
		return
	}
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	tl := new(tunnel)
	tl.local = c.LocalAddr().String()
	tl.remote = c.RemoteAddr().String()
	log.Println("local:", c.LocalAddr(),"remote:", c.RemoteAddr())
	log.Println("network local:", c.LocalAddr().Network(), c.RemoteAddr().Network())

	var b [1024]byte
	n, err := c.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	//var method, host, address string
	sscanf, err := fmt.Sscanf(string(b[:bytes.IndexByte(b[:],'\n')]), "%s%s", &tl.method, &tl.host)
	if err != nil {
		log.Println(sscanf, err)
		return
	}

	log.Println("method:", tl.method, "host:", tl.host)
	hostURL, err := url.Parse(tl.host)
	if err != nil {
		return
	}

	log.Println("host:", tl.host,
		"Host:", hostURL.Host,
		"Path:", hostURL.Path,
		"Scheme:", hostURL.Scheme,
		"Opaque:", hostURL.Opaque)

	tl.host = hostURL.Host
	tl.hostName = hostURL.Hostname()
	tl.hostPort = hostURL.Port()

	if forwardPort == "" {
		if tl.ListenPort == "" {
			tl.dialAddr = tl.hostName + ":80"
		}
	} else {
		forwardPort = strings.Trim(forwardPort,": ")
		tl.dialAddr = tl.hostName + ":" + forwardPort
	}

	if hostURL.Opaque == "443" {
		tl.dialAddr = hostURL.Scheme + ":443"
	} else {
		if strings.Index(hostURL.Host, ":") == -1 {
		} else {
			tl.dialAddr = hostURL.Host
		}
	}

	server, err := net.Dial("tcp", tl.dialAddr)
	log.Println("Dial server Addr:", tl.dialAddr)
	if err != nil {
		log.Println(err)
		return
	}
	if tl.method == "CONNECT" || tl.method == "connect" {
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
