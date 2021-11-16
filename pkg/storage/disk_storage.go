package storage

import (
	"k8s.io/klog/v2"
	"os"
)

const (
	CacheBaseDir = "/etc/kubernetes/cache/"
)

type cacheStorage struct {
	baseDir          string
}

// CreateStorage create a storage.Store for backend storage
func CreateStorage(cachePath string) (*cacheStorage, error) {
	if cachePath == "" {
		klog.Infof("disk cache path is empty, set it by default %s", CacheBaseDir)
		cachePath = CacheBaseDir
	}
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		if err = mkdir(cachePath); err != nil {
			return nil, err
		}
	}

	ds := &cacheStorage{
		baseDir:	cachePath,
	}

	return ds, nil
}

func mkdir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		klog.Fatalf("mkdir %s error: %v", dirPath, err)
		return err
	}
	return nil
}