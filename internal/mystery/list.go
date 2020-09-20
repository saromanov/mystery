package mystery

import "fmt"

// ListRequest defines struct for getting request
type ListRequest struct {
	MasterPass string
}

// ListResponse defines struct for getting response
type ListResponse struct {
}

// List returns list of secrets
func List(p ListRequest) (Data, error) {
	rsp, err := p.Backend.List([]byte(p.MasterPass))
	if err != nil {
		return "", fmt.Errorf("list: unable to get data: %v", err)
	}
	return ListResponse{}, nil
}
