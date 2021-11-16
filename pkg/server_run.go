package main

import (
	"fmt"
	_ "github.com/spf13/cobra"
	"github.com/swiftrivergo/snedge/pkg/proxy"
	"github.com/swiftrivergo/snedge/pkg/storage"
	boltstorage "github.com/swiftrivergo/snedge/pkg/storage/bolt"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"os"
)

const target = "172.16.0.5:8081"
const source = "127.0.0.1:8081"
const protocol = "http://"

func main() {

	//NewBoltStorage
	//default dbfile: CacheBaseDir = "/etc/kubernetes/cache/"

	s, err := storage.CreateStorage("")
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("storage dir:", s.CacheStorageDir())
	store, err := boltstorage.NewBoltStorage(s.CacheStorageDir())
	if err != nil {
		klog.Errorln(err)
	}

	key := "hello"
	value := "bolt"
	err = store.Create(key, []byte(value))
	if err != nil {
		klog.Errorln(err)
	}

	data, err := store.Get(key)
	if err != nil {
		klog.Errorln(err)
	}
	fmt.Println("hello:", string(data[:]))

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
