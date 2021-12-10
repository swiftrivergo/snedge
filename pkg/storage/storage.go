package storage

type Storage interface {
	Create(key string, cache []byte) error
	Delete(key string) error
	Update(key string, cache []byte) error
	Get(key string) ([]byte, error)
	Path() string
}
