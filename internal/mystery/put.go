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
	Data       Data
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

	value, err := p.Data.encode()
	if err != nil {
		return fmt.Errorf("unable to encode value: %v", err)
	}
	reducedValue, err := compress(value)
	if err != nil {
		return fmt.Errorf("unable to compress data: %v", err)
	}
	compressed := len(reducedValue) < len(value)
	if compressed {
		value = reducedValue
	}
	if err := p.Backend.Put([]byte(p.MasterPass), backend.Secret{
		Namespace:  []byte(p.Namespace),
		Data:       value,
		Compressed: compressed,
	}); err != nil {
		return fmt.Errorf("put: unable to store data: %v", err)
	}
	return nil
}
