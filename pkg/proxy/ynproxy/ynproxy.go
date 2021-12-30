package ynproxy

import (
	"fmt"
	"github.com/swiftrivergo/snedge/pkg/util"
	"io"
	"io/ioutil"
	_ "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/klog/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var target *string

type director func(r *http.Request)

type YNProxy struct {
	target *url.URL
	director director
}

type Proxy interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
}

func buildDirector(r *http.Request) director {
	d := func(req *http.Request) {
		req = r
	}

	return d
}

func NewYNProxy() *YNProxy {
	s := "http://127.0.0.1"
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return &YNProxy{
		target:   u,
		director: nil,
	}
}

func NewDestProxy(target *url.URL) *YNProxy {
	if target == nil {
		return nil
	}

	p := NewYNProxy()
	p.SetTarget(target)

	fmt.Println("p:", p.target.Scheme,p.target.Host,p.target.Path)
	return p
}

func (p *YNProxy) Target() *url.URL {
	return p.target
}

func (p *YNProxy) SetTarget(url *url.URL) {
	p.target = url
}

func (p *YNProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	serveHTTP(rw, req)
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {

	uri := *target + r.RequestURI
	klog.Infoln(r.Method + ": " + uri)

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			klog.Errorln(err)
		}
		klog.Infoln("Body: %v\n", string(body))
	}

	rr, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		klog.Errorln(err)
	}

	copyHeader(r.Header, &rr.Header)

	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	klog.Infoln(err)
	klog.Infoln("x-forwarded-for: %v\n", resp.Header)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			klog.Errorln(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Errorln(err)
	}

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("x-forwarded-for", rr.Host)

	_, err = w.Write(body)
	if err != nil {
		klog.Errorln(err)
		return
	}
}

func copyHeader(source http.Header, dest *http.Header){
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func (p *YNProxy) direct(r *http.Request) *http.Request {
	director := buildDirector(r)
	p.director = director
	return customDirectorPolicy(p.target, r)
}

func customDirectorPolicy(target *url.URL, r *http.Request) *http.Request {
	r.URL.Scheme = "http"
	r.URL.Host = target.Host
	r.URL.Path = target.Path
	return r
}

func (p *YNProxy) ReverseProxy(target *url.URL) *httputil.ReverseProxy {
	p.target = target
	return util.NewReverseProxy(target)
}
