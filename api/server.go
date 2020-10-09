package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/saromanov/mystery/config"
	"github.com/saromanov/mystery/internal/mystery"
)

// Make provides starting of the server
func Make(c *config.Server, mys *mystery.Mystery) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	a := API{}
	r.Get("/", a.Put)
	if err := http.ListenAndServe(c.Address, r); err != nil {
		return fmt.Errorf("unable to init server: %v", err)
	}
	return nil
}
