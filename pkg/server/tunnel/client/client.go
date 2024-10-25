package client

import (
	"context"
	"fmt"
	"github.com/swiftrivergo/snedge/pkg/server/tunnel/transport"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/rancher/remotedialer"
	"github.com/twinj/uuid"
)

// Client for ynCmd
type Client struct {
	// Remote site for websocket address
	Remote string

	// Map of upstream servers dns.entry=http://ip:port
	UpstreamMap map[string]string

	// Token for authentication
	Token string

	// StrictForwarding
	StrictForwarding bool
}

func makeAllowsAllFilter() func(network, address string) bool {
	return func(network, address string) bool {
		return true
	}
}

func makeFilter(upstreamMap map[string]string) func(network, address string) bool {
	trimmedMap := map[string]bool{}

	for _, v := range upstreamMap {
		u, err := url.Parse(v)
		if err != nil {
			log.Printf("Error parsing: %s, skipping.\n", v)
			continue
		}

		trimmedMap[u.Host] = true
	}

	return func(network, address string) bool {
		if network != "tcp" {
			log.Printf("network not allowed: %q\n", network)

			return false
		}

		if ok, v := trimmedMap[address]; ok && v {
			return true
		}

		return false
	}
}

// Connect and serve traffic through websocket
func (c *Client) Connect() error {
	headers := http.Header{}
	headers.Set(transport.InHeader, uuid.Formatter(uuid.NewV4(), uuid.FormatHex))
	for k, v := range c.UpstreamMap {
		headers.Add(transport.UpstreamHeader, fmt.Sprintf("%s=%s", k, v))
	}
	if c.Token != "" {
		headers.Add("Authorization", "Bearer "+c.Token)
	}

	remote := c.Remote
	if !strings.HasPrefix(remote, "ws") {
		remote = "ws://" + remote
	}
	var filter func(network, address string) bool

	if c.StrictForwarding {
		filter = makeFilter(c.UpstreamMap)
	} else {
		filter = makeAllowsAllFilter()
	}

	for {
		remotedialer.ClientConnect(context.Background(), remote+ "/tunnel", headers, nil, filter, nil)
	}
}
