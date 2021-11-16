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

var base = ""

type cacheStorage struct {
	basePath string
}

func (c cacheStorage) StorageBasePath() string {
	return c.basePath
}

func DefaultStorageBasePath() (string, error) {
	system := runtime.GOOS
	if system == "linux" {
		base = CacheBaseLinuxDir
	} else if system == "windows" {
		base = CacheBaseWindowDir
	}

	if _, err := os.Stat(filepath.Dir(base)); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return base, err
		}
	} else {
		return base, nil
	}
}

// CreateStorage create a storage.Store for backend storage
func CreateStorage(cachePath string) (*cacheStorage, error) {
	if cachePath == "" {
		system := runtime.GOOS
		if system == "linux" {
			base = CacheBaseLinuxDir
		} else if system == "windows" {
			base = CacheBaseWindowDir
		}
		klog.Infof("disk cache path is empty, set it by default %s", base)
	} else {
		base = cachePath
	}
	if _, err := os.Stat(filepath.Dir(base)); os.IsNotExist(err) {
		if err = mkdir(filepath.Dir(base)); err != nil {
			return nil, err
		}
	}

	store := &cacheStorage{
		basePath: base,
	}

	return store, nil
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