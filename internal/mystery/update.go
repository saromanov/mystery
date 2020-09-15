package mystery

import (
	"fmt"

	"github.com/saromanov/mystery/internal/backend"
)

// UpdateRequest provides struct for sending to update
type UpdateRequest struct {
	MasterPass string
	Namespace  string
	Backend    backend.Backend
	Data       Data
}

// validate provides validation of data
func (p UpdateRequest) validate() error {
	if p.MasterPass == "" {
		return errNoMasterPass
	}
	if p.Namespace == "" {
		return errNoNamespace
	}
	if p.Backend == nil {
		return errNoBackend
	}
	return nil
}

// Update provides updating value by the key
func Update(p UpdateRequest) error {
	if err := p.validate(); err != nil {
		return fmt.Errorf("update: unable to validate data: %v", err)
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
	fmt.Println("DATA: ", value)
	err = p.Backend.Update([]byte(p.MasterPass), backend.Secret{
		Namespace:  []byte(p.Namespace),
		Compressed: compressed,
		Data:       value,
	})
	if err != nil {
		return fmt.Errorf("update: unable to get data: %v", err)
	}
	return nil
}
