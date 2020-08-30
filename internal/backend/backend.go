package backend

// Backend defines way for store secrets
type Backend interface {
	Get() ([]byte, error)
	Put([]byte) error
}

// MasterPassBackend defines backend for master pass
type MasterPassBackend interface {
}
