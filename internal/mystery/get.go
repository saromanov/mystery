package mystery

import (
	"fmt"

	"github.com/saromanov/mystery/internal/backend"
)

// GetRequest provides struct for sending to Get
type GetRequest struct {
	MasterPass string
	Namespace  string
	Backend    backend.Backend
}

// validate provides validation of data
func (p GetRequest) validate() error {
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

// Get provides getting value by the key
func Get(p GetRequest) (Data, error) {
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("get: unable to validate data: %v", err)
	}
	rsp, err := p.Backend.Get([]byte(p.MasterPass), []byte(p.Namespace))
	if err != nil {
		return nil, fmt.Errorf("get: unable to get data: %v", err)
	}
	return Decode(rsp.Data)
}
