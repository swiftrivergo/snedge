package storage

import (
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
	"runtime"
)

const (
	CacheBaseLinuxDir = "/etc/kubernetes/cache"
	CacheBaseWindowDir = "c:/etc/kubernetes/cache"
)

var basedir = ""

type cacheStorage struct {
	baseDir          string
}

func (c cacheStorage) CacheStorageDir() string {
	return basedir
}

// CreateStorage create a storage.Store for backend storage
func CreateStorage(cachePath string) (*cacheStorage, error) {
	if cachePath == "" {
		system := runtime.GOOS
		if system == "linux" {
			basedir = CacheBaseLinuxDir
		} else if system == "windows" {
			basedir = CacheBaseWindowDir
		}
		klog.Infof("disk cache path is empty, set it by default %s", basedir)
	}
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		if err = mkdir(filepath.Dir(basedir)); err != nil {
			return nil, err
		}
	}

	ds := &cacheStorage{
		baseDir:	basedir,
	}

	return ds, nil
}

func mkdir(dirPath string) error {
	fmt.Println("dirPath:", dirPath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		klog.Fatalf("mkdir %s error: %v", dirPath, err)
		return err
	}
	return nil
}