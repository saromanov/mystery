package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/saromanov/mystery/config"
	"github.com/saromanov/mystery/internal/backend"
)

// Postgres defines backend for postgres
type Postgres struct {
	db *gorm.DB
}

// New provides initialization of postgres for store master pass
func New(c *config.MasterPassBackend) (backend.Backend, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Host, c.Port, c.User, c.DB, c.Password)
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	return &Postgres{
		db: db,
	}, nil
}

// Get defines getting a secret from backend
func (m *Postgres) Get(key string) ([]byte, error) {
	return nil, nil
}

// Put defines putting a secret to backend
func (m *Postgres) Put(key string) error {
	return nil
}
