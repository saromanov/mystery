package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/saromanov/mystery/config"
	"github.com/saromanov/mystery/internal/mystery"
	log "github.com/sirupsen/logrus"
)

// Make provides starting of the server
func Make(c *config.Server, l *log.Logger, mys *mystery.Mystery) error {
	if c == nil {
		return fmt.Errorf("config is not defined")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	a := New(mys)
	r.Get("/put", a.Put)
	l.Infof("starting of server at address %s...", c.Address)
	if err := startServer(c, r); err != nil {
		return fmt.Errorf("unable to init server: %v", err)
	}
	return nil
}

func startServer(c *config.Server, h http.Handler) error {
	if c.Key != "" && c.Crt != "" {
		return startTLSServer(c, h)
	}

	if err := http.ListenAndServe(c.Address, h); err != nil {
		return fmt.Errorf("unable to init server: %v", err)
	}
	return nil
}

func startTLSServer(c *config.Server, h http.Handler) error {
	err := http.ListenAndServeTLS(c.Address, c.Crt, c.Key, h)
	if err != nil {
		return fmt.Errorf("unable to start tls server: %v", err)
	}
	return nil
}
