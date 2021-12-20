package proxy

import (
	"errors"
	"fmt"
	_ "k8s.io/apimachinery/pkg/util/rand"
	"math/rand"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	if target == nil {
		fmt.Println(target)
		panic(any(errors.New("<nil>")))
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func NewRandReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	i := rand.Int()%len(targets)
	target := targets[rand.Int()%len(targets)]
	fmt.Println("rand url index:", i, target)
	return NewReverseProxy(target)
}

func NewReverseProxies(targets []*url.URL) []*httputil.ReverseProxy {
	fmt.Println(len(targets))
	sfr := make([]*httputil.ReverseProxy,0)
	for _, v := range targets {
		value := v
		fmt.Println("url:", value.Scheme, value.Host, value.Path,
			",", value.Hostname(),value.RequestURI(),value.Port())
		fr := NewReverseProxy(value)
		sfr = append(sfr, fr)
	}
	return sfr
}
