package stores

// Storage defines the requireed interface of the storagee used by stores defined in this package.
type Storage interface {
	UpdateInts(f func(vals []int) error, keys ...string) ([]int, error)
	GetInt(key string, default_ *int) (int, error)
	SetInt(key string, x int) error

	UpdateJSON(key string, v interface{}, f func() error) error
	GetJSON(key string, v interface{}) (bool, error)
}
