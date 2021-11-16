package proxy

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var target *string

type Proxy struct {
	target *url.URL
}

func NewProxy() *Proxy {
	target = flag.String("target", "http://127.0.0.1", "target URL for reverse proxy")
	flag.Parse()

	parse, err := url.Parse(*target)
	if err != nil {
		return nil
	}
	return &Proxy{target: parse}
}

func (p *Proxy) Target() *url.URL {
	return p.target
}

func (p *Proxy) SetTarget(url *url.URL) {
	p.target = url
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	serveHTTP(rw, p.rebuildRequestHost(req))
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
	klog.Infoln("Resp-Headers: %v\n", resp.Header)

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
	dH.Add("Requested-Host", rr.Host)

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

func (p *Proxy) ReverseProxy(target *url.URL) *httputil.ReverseProxy {
	p.target = target
	return ReverseProxy(target)
}

func (p *Proxy) rebuildRequestHost(r *http.Request) *http.Request {
	return rebuildRequestHost(p.target, r)
}

func ReverseProxy(target *url.URL) *httputil.ReverseProxy {
	if target == nil {
		panic(errors.New("target is nil"))
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func rebuildRequestHost(target *url.URL, r *http.Request) *http.Request {
	r.URL.Host = target.Host
	r.Host = target.Host
	return r
}