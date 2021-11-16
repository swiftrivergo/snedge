package main

import (
	_ "github.com/spf13/cobra"
	"github.com/swiftrivergo/snedge/pkg/proxy"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
)

const target = "172.16.0.5:8081"
const source = "127.0.0.1:8081"
const protocol = "http://"

func main() {

	if dest, err := url.Parse(protocol+target); err != nil {
		klog.Errorln(err)
	} else {
		p := proxy.NewProxy()
		p.SetTarget(dest)
		if err := http.ListenAndServe(source, proxy.ReverseProxy(dest)); err != nil {
			klog.Errorln(err)
		}
	}
}
