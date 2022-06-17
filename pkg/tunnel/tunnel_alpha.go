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

type tunnel struct {
	tcp  string
	//Addr string
	addr string
	forwardPort string

	listenAddr string
	listenPort string

	dialAddr string
	dialPort string

	method string
	host string
	hostName string
	hostPort string

	local string
	remote string
}

var tl *tunnel

func init() {
	tl = &tunnel{
		tcp:         "",
		addr:        "",
		forwardPort: "",
		listenAddr:  "",
		listenPort:  "",
		dialAddr:    "",
		dialPort:    "",
		method:      "",
		host:        "",
		hostName:    "",
		hostPort:    "",
		local:       "",
		remote:      "",
	}
}

func New() *tunnel {
	tl = new(tunnel)
	//The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	tl.tcp = "tcp"
	return tl
}

func (t *tunnel) GetListenAddr() string {
	return t.listenAddr + ":" + t.listenPort
}

func (t *tunnel) BindForwardPort(port string) {
	t.forwardPort = port
}

func (t *tunnel) bindListenAddr(address string) {
	t.addr = address
	addr := strings.Trim(address, " ")
	if strings.Index(addr, ":") == -1 {
		t.listenAddr = ":80"
	} else {
		t.listenAddr = strings.Split(addr, ":")[0]
		t.listenPort = strings.SplitAfter(strings.Trim(addr, " "), ":")[1]
		log.Println(address, addr, t.tcp, t.listenAddr, t.listenPort)
	}
}

func (t *tunnel) AddListenAddr(address string) {
	t.bindListenAddr(address)
}

func (t *tunnel) Run() error {
	l, err := net.Listen(t.tcp, t.GetListenAddr())

	log.Println("listen tcp:", t.tcp, "GetListenAddr():", t.GetListenAddr())
	if err != nil {
		log.Panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleConnForward(c, t.forwardPort)
	}
	return nil
}

func handleConnForward(c net.Conn, forwardPort string) {
	var b [1024]byte

	if c == nil {
		return
	}
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	tl.local = c.LocalAddr().String()
	tl.remote = c.RemoteAddr().String()

	n, err := c.Read(b[:])
	if err != nil {
		return
	}

	sscanf, err := fmt.Sscanf(string(b[:bytes.IndexByte(b[:],'\n')]), "%s%s", &tl.method, &tl.host)
	if err != nil {
		log.Println(sscanf, err)
		return
	}

	hostURL, err := url.Parse(tl.host)
	if err != nil {
		return
	}

	tl.host = hostURL.Host
	tl.hostName = hostURL.Hostname()
	tl.hostPort = hostURL.Port()

	if forwardPort == "" {
		if tl.listenPort == "" {
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

	log.Println("Dial server Addr:", tl.dialAddr)
	server, err := net.Dial("tcp", tl.dialAddr)
	if err != nil {
		return
	}
	if tl.method == "CONNECT" || tl.method == "connect" {
		fprint, err := fmt.Fprint(c, "HTTP/1.1 200 Connection establish\r\n")
		if err != nil {
			log.Println(fprint, err)
			return
		}
	} else {
		_, err := server.Write(b[:n])
		if err != nil {
			return
		}
	}

	go func() {
		_, err := io.Copy(server, c)
		if err != nil {
		}
	}()
	_, err = io.Copy(c, server)
	if err != nil {
		return
	}
}
