package bolt

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/swiftrivergo/snedge/pkg/storage"
	bolt "go.etcd.io/bbolt"
	"k8s.io/klog/v2"
	"path/filepath"
	_ "path/filepath"
	"strings"
	"sync"
	"time"
)

const Bucket = "XxxEdge"

type boltStorage struct {
	basePath string
	db *bolt.DB
	sync.Mutex
}

func NewBoltStorage(dbFile string) (storage.Storage, error) {
	fmt.Println("db dir:", filepath.Dir(dbFile))
	_, err := storage.CreateStorage(dbFile)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(dbFile)
	if base == "" {
		return nil, errors.New("file is invalid")
	}

	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		klog.Fatalf("init bolt storage error: %v", err)
		return nil, err
	}

	onstorage := &boltStorage{db: db}
	onstorage.basePath = dbFile

	// init
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		klog.Fatalf("init bolt storage error: %v", err)
		return onstorage, err
	}

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists([]byte(Bucket))
	if err != nil {
		klog.Fatalf("init bolt storage error: %v", err)
		return onstorage, nil
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		klog.Fatalf("init bolt storage error: %v", err)
		return onstorage, nil
	}

	return onstorage, nil
}

func (bs *boltStorage) Create(key string, data []byte) error {
	bs.Lock()
	defer bs.Unlock()
	return bs.create(key, data)
}

func (bs *boltStorage) create(key string, data []byte) error {
	klog.V(8).Infof("storage one key=%s, cache=%s", key, string(data))

	err := bs.update(bs.oneKey(key), data)
	if err != nil {
		klog.Errorf("write one cache %s error: %v", key, err)
		return err
	}

	return nil
}

func (bs *boltStorage) Update(key string, cache []byte) error {
	bs.Lock()
	defer bs.Unlock()
	return bs.update(key, cache)
}

func (bs *boltStorage) update(key string, data []byte) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Bucket))
		err := b.Put([]byte(key), data)
		return err
	})
	return err
}

func (bs *boltStorage) Get(key string) ([]byte, error) {
	return bs.get(key)
}

func (bs *boltStorage) get(key string) ([]byte, error) {
	var data []byte
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Bucket))
		data = b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("no data for %s", key)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bs *boltStorage) Delete(key string) error {
	bs.Lock()
	defer bs.Unlock()
	return bs.delete(key)
}

func (bs *boltStorage) delete(key string) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Bucket))
		err := b.Delete([]byte(key))
		return err
	})
	return err
}

func (bs *boltStorage) storeList(key string, data []byte) error {
	klog.V(8).Infof("storage list key=%s, cache=%s", key, string(data))

	err := bs.update(bs.listKey(key), data)
	if err != nil {
		klog.Errorf("write list cache %s error: %v", key, err)
		return err
	}

	return nil
}

func (bs *boltStorage) loadOne(key string) ([]byte, error) {
	data, err := bs.get(bs.oneKey(key))
	if err != nil {
		klog.Errorf("read one cache %s error: %v", key, err)
		return nil, err
	}

	klog.V(8).Infof("load one key=%s, cache=%s", key, string(data))
	return data, nil
}

func (bs *boltStorage) loadList(key string) ([]byte, error) {
	data, err := bs.get(bs.listKey(key))
	if err != nil {
		klog.Errorf("read list cache %s error: %v", key, err)
		return nil, err
	}

	klog.V(8).Infof("load list key=%s, cache=%s", key, string(data))
	return data, nil
}

func (bs *boltStorage) oneKey(key string) string {
	return strings.ReplaceAll(key, "/", "_")
}

func (bs *boltStorage) listKey(key string) string {
	return strings.ReplaceAll(fmt.Sprintf("%s_%s", key, "list"), "/", "_")
}

func (bs *boltStorage) Path() string {
	return bs.basePath
}