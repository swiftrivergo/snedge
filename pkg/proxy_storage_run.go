package main

import (
	"fmt"
	_ "github.com/spf13/cobra"
	"github.com/swiftrivergo/snedge/pkg/proxy/ynproxy"
	"github.com/swiftrivergo/snedge/pkg/storage"
	boltstorage "github.com/swiftrivergo/snedge/pkg/storage/bolt"
	"github.com/swiftrivergo/snedge/pkg/tunnel"
	"k8s.io/klog/v2"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	target = "172.16.0.5:8081"
	target1 = "127.0.0.1:8081"
	target2 = "127.0.0.2:8081"
	target3 = "127.0.0.3:8081"
	target4 = "127.0.0.4:8081"
	target5 = "127.0.0.5:8081"
	target6 = "127.0.0.6:8081"
	target7 = "127.0.0.7:8081"
	target8 = "127.0.0.8:8081"
	target9 = "127.0.0.9:8081"
	target10 = "127.0.0.10:8081"
	source = "127.0.0.1:8081"
	protocol = "http://"
)

func main() {

	//NewBoltStorage
	//default dbfile: CacheBaseDir = "/etc/kubernetes/cache/"
	rand.Seed(time.Now().UnixNano())

	s, err := storage.CreateStorage("")
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("storage:", s.StorageBasePath())
	store, err := boltstorage.NewBoltStorage(s.StorageBasePath())
	if err != nil {
		klog.Errorln(err)
	}
	fmt.Println("bolt storage path:", store.Path())
	baseBath, _ := storage.DefaultStorageBasePath()
	fmt.Println("default storage base path:", baseBath)

	key := "hello"
	err = store.Create(key, []byte(""))
	if err != nil {
		klog.Errorln(err)
	}

	value := strconv.Itoa(rand.Int())
	fmt.Println("rand value:", value)
	err = store.Update(key, []byte(value))
	if err != nil {
		klog.Errorln(err)
	}

	data, err := store.Get(key)
	if err != nil {
		klog.Errorln(err)
	}
	fmt.Println("get  value:", string(data[:]))

	if dest, err := url.Parse(protocol+target1); err != nil {
		klog.Errorln(err)
	} else {
		urls := make([]*url.URL,0)
		urls = append(urls, dest)

		if url2, err := url.Parse(protocol+target2); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url2)
		}

		if url3, err := url.Parse(protocol+target3); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url3)
		}

		if url4, err := url.Parse(protocol+target4); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url4)
		}

		if url5, err := url.Parse(protocol+target5); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url5)
		}

		if url6, err := url.Parse(protocol+target6); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url6)
		}

		if url7, err := url.Parse(protocol+target7); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url7)
		}

		if url8, err := url.Parse(protocol+target8); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url8)
		}

		if url9, err := url.Parse(protocol+target9); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url9)
		}

		if url10, err := url.Parse(protocol+target10); err != nil {
			fmt.Println(err)
		} else {
			urls = append(urls, url10)
		}

		//v1:
		//proxies := proxy.NewReverseProxies(urls)
		//fmt.Println(proxies)
		//e := funcListenAndServe(proxies)
		//if e != nil {
		//	fmt.Println(e)
		//}

		//v2:
		tl := tunnel.New()
		tl.Addr = ":8082"

		go func() {
			//Todo: port should be set by user
			err := tl.Listen()
			if err != nil {
				log.Println(err)
			}
		}()

		select {
		}
	}
}

func ListenAndServe(proxies []*ynproxy.YNProxy) error {
	var err error
	for i, v := range proxies {
		value := v
		port := fmt.Sprintf(":%d", 8082 + i)
		fmt.Println("port:", port)

		go func(v *ynproxy.YNProxy) {
			fmt.Println("target host:",value,",port:",port)
			e := http.ListenAndServe(port, v)
			if e != nil {
				err = e
				fmt.Println(err)
			}
		}(value)
	}

	return err
}


func funcListenAndServe(proxies []*httputil.ReverseProxy) error {
	var err error
	for i, v := range proxies {
		val := make(chan *httputil.ReverseProxy, 1)
		val <- v
		port := fmt.Sprintf(":%d", 8082 + i)
		fmt.Println("port:", port)

		go func(rp chan *httputil.ReverseProxy, port string ) {
			vv := <- rp
			vp := port
			fmt.Println("target host:", vv, ", port:", vp)
			e := http.ListenAndServe(vp, vv)
			if e != nil {
				err = e
				fmt.Println(err)
			}
		}(val, port)
	}

	return err
}