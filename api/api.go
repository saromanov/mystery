package api

import (
	"encoding/json"
	"net/http"

	"github.com/saromanov/mystery/internal/mystery"
	"github.com/saromanov/mystery/internal/backend"
)

// API defines configuration for restfuil api
type API struct {
	core *mystery.Mystery
	masterPass string
	backend backend.Backend

}

// New provides initialization of API
func New(m *mystery.Mystery) *API {
	return &API{
		core: m,
	}
}

func (a *API) Put(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    var m PutRequest
    err := decoder.Decode(&m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	data, storeType := mystery.Prepare(m.Value)
	if err := a.core.Put(mystery.PutRequest{
		Namespace: m.Namespace,
		Data: data,
		Type: storeType, 
		Backend: a.backend,
		MasterPass: a.masterPass,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get provides getting of getting data
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	lst, err := a.core.Get(mystery.GetRequest{
		Namespace: namespace,
	})
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
