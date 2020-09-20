package mystery

import (
	"fmt"
	"time"

	"github.com/saromanov/mystery/internal/backend"
)

// ListRequest defines struct for getting request
type ListRequest struct {
	MasterPass string
	Backend    backend.Backend
}

// ListResponse defines struct for getting response
type ListItemResponse struct {
	Namespace      string
	Data           []byte
	UserID         string
	CreatedAt      time.Time
	CurrentVersion uint64
	MaxVersion     uint64
}

// List returns list of secrets
func List(p ListRequest) ([]ListItemResponse, error) {
	rsp, err := p.Backend.List([]byte(p.MasterPass))
	if err != nil {
		return nil, fmt.Errorf("list: unable to get data: %v", err)
	}
	data := make([]ListItemResponse, len(rsp))
	for i, r := range rsp {
		data[i] = ListItemResponse{
			Namespace:      r.Namespace,
			Data:           r.Data,
			MaxVersion:     r.MaxVersion,
			CurrentVersion: r.CurrentVersion,
			UserID:         r.UserID,
			CreatedAt:      r.CreatedAt,
		}
	}
	return data, nil
}
