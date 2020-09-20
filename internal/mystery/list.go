package mystery

import (
	"fmt"

	"github.com/saromanov/mystery/internal/backend"
)

// ListRequest defines struct for getting request
type ListRequest struct {
	MasterPass string
	Backend    backend.Backend
}

// ListResponse defines struct for getting response
type ListResponse struct {
}

// List returns list of secrets
func List(p ListRequest) (ListResponse, error) {
	_, err := p.Backend.List([]byte(p.MasterPass))
	if err != nil {
		return ListResponse{}, fmt.Errorf("list: unable to get data: %v", err)
	}
	return ListResponse{}, nil
}
