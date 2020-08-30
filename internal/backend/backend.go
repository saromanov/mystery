package backend

// Backend defines way for store secrets
type Backend interface {
	Get(masterKey, key []byte) ([]byte, error)
	Put(masterKey, key, value []byte) error
}

// MasterPassBackend defines backend for master pass
type MasterPassBackend interface {
	Get(pass string) ([]byte, error)
	Put(pass string) error
}
