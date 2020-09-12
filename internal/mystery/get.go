package mystery

import (
	"bytes"
	"encoding/gob"
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
func Get(p GetRequest) (map[string]string, error) {
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("get: unable to validate data: %v", err)
	}
	rsp, err := p.Backend.Get([]byte(p.MasterPass), []byte(p.Namespace))
	if err != nil {
		return nil, fmt.Errorf("get: unable to get data: %v", err)
	}
	return decodeValue(rsp.Data)
}

func decodeValue(data []byte) (map[string]string, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var out map[string]string
	if err := dec.Decode(&out); err != nil {
		return out, fmt.Errorf("unable to decode value: %v", err)
	}
	return out, nil
}
