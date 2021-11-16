package main

import (
	_ "github.com/spf13/cobra"
	"github.com/swiftrivergo/snedge/pkg/proxy"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
)

const target = "http://172.16.0.5:8081"

func main() {

	if dest, err := url.Parse(target); err != nil {
		klog.Errorln(err)
	} else {
		if err := http.ListenAndServe("127.0.0.1:8081", proxy.ReverseProxy(dest)); err != nil {
			klog.Errorln(err)
		}
	}
}
