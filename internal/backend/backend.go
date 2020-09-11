package backend

import "time"

// Backend defines way for store secrets
type Backend interface {
	Get(masterKey, key []byte) (Secret, error)
	Put(masterKey []byte, secret Secret) error
}

// MasterPassBackend defines backend for master pass
type MasterPassBackend interface {
	Get(pass string) ([]byte, error)
	Put(pass string) error
}

// Secret defines struct for store secrets
type Secret struct {
	Namespace    []byte
	Data         []byte
	ExpiredAfter *time.Duration
}
