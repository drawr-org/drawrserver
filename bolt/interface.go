package bolt

// DBClient is the interface that gets implemented by this package
type DBClient interface {
	Open() error
	Close()
	Get(bucket, key string) ([]byte, error)
	Put(bucket, key string, value interface{}) error
}
