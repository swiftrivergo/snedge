package storage

type Storage interface {
	Create(key string, cache []byte) error
	Delete(key string) error
	Update(key string, cache []byte) error
	Get(key string) ([]byte, error)

	// GetOne Get One
	//GetOne() (string, []byte)

	//Exist(key string) (bool, error)
	//List(key string) ([][]byte, error)
	//ListKeys(key string) ([]string, error)
}
