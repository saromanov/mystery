package backend

import "time"

// Backend defines way for store secrets
type Backend interface {
	Get(masterKey, key []byte) (Secret, error)
	Put(masterKey []byte, secret Secret) error
	Update(masterKey []byte, secret Secret) error
	Delete(masterKey, key []byte) error
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
	Compressed   bool
}

// DeleteSecret defines request for delete secret
type DeleteSecret struct {
	Namespace []byte
	Version   int
	Force     bool
}
