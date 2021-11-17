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
const target2 = "127.0.0.2:8081"
const target3 = "127.0.0.3:8081"
const source = "127.0.0.1:8081"
const protocol = "http://"

func main() {

	//NewBoltStorage
	//default dbfile: CacheBaseDir = "/etc/kubernetes/cache/"

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
	//value := "bolt"
	//err = store.Create(key, []byte(value))
	//if err != nil {
	//	klog.Errorln(err)
	//}

	data, err := store.Get(key)
	if err != nil {
		klog.Errorln(err)
	}
	fmt.Println("hello:", string(data[:]))

	if dest, err := url.Parse(protocol+target); err != nil {
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

		randProxy := proxy.NewRandReverseProxy(urls)

		fmt.Println("urls:", urls)
		err := http.ListenAndServe(":8082", randProxy)
		if err != nil {
			fmt.Println(err)
		}
	}

}
