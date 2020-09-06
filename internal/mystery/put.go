package mystery

import (
	"errors"
	"fmt"

	"github.com/saromanov/mystery/internal/backend"
)

var (
	errNoMasterPass = errors.New("master pass is not defined")
	errNoKey        = errors.New("key is not defined")
	errNoValue      = errors.New("value is not defined")
	errNoBackend    = errors.New("backend is not defined")
)

// PutRequest provides struct for sending to Put
type PutRequest struct {
	MasterPass string
	Key        string
	Value      string
	Backend    backend.Backend
}

// validate provides validation of data
func (p PutRequest) validate() error {
	if p.MasterPass == "" {
		return errNoMasterPass
	}
	if p.Key == "" {
		return errNoKey
	}
	if p.Value == "" {
		return errNoValue
	}
	if p.Backend == nil {
		return errNoBackend
	}
	return nil
}

// Put provides adding key-value pair to backend
func Put(p PutRequest) error {
	if err := p.validate(); err != nil {
		return fmt.Errorf("put: unable to validate data: %v", err)
	}
	if err := p.Backend.Put([]byte(p.MasterPass), backend.Secret{
		Key:   []byte(p.Key),
		Value: []byte(p.Value),
	}); err != nil {
		return fmt.Errorf("put: unable to store data: %v", err)
	}
	return nil
}
