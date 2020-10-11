package api

import (
	"encoding/json"
	"net/http"

	"github.com/saromanov/mystery/internal/mystery"
)

// API defines configuration for restfuil api
type API struct {
	core *mystery.Mystery
}

// New provides initialization of API
func New(m *mystery.Mystery) *API {
	return &API{
		core: m,
	}
}

func (a *API) Put(w http.ResponseWriter, r *http.Request) {
	if err := a.core.Put(mystery.PutRequest{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get provides getting of
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	lst, err := a.core.Get(mystery.GetRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(lst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
