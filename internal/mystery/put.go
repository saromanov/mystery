package mystery

import (
	"errors"
	"fmt"

	"github.com/saromanov/mystery/internal/backend"
)

var (
	errNoMasterPass = errors.New("master pass is not defined")
	errNoNamespace  = errors.New("namespace is not defined")
	errNoValue      = errors.New("value is not defined")
	errNoBackend    = errors.New("backend is not defined")
)

// PutRequest provides struct for sending to Put
type PutRequest struct {
	MasterPass string
	Namespace  string
	Data       map[string]string
	Backend    backend.Backend
}

// validate provides validation of data
func (p PutRequest) validate() error {
	if p.MasterPass == "" {
		return errNoMasterPass
	}
	if p.Namespace == "" {
		return errNoNamespace
	}
	if len(p.Data) == 0 {
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
