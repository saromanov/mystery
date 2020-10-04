package mystery

import (
	"fmt"

	"github.com/saromanov/mystery/internal/backend"
)

// DeleteRequest provides struct for sending to Delete
type DeleteRequest struct {
	MasterPass string
	Namespace  string
	Backend    backend.Backend
	Version    int64
}

// validate provides validation of data
func (p DeleteRequest) validate() error {
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

// Delete provides getting value by the key
func (m *Mystery) Delete(p DeleteRequest) error {
	if err := p.validate(); err != nil {
		return fmt.Errorf("delete: unable to validate data: %v", err)
	}
	err := p.Backend.Delete([]byte(p.MasterPass), []byte(p.Namespace))
	if err != nil {
		return fmt.Errorf("delete: unable to get data: %v", err)
	}

	return nil
}
