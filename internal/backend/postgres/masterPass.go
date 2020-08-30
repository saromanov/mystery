package postgers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/saromanov/mystery/config"
	"github.com/saromanov/mystery/internal/backend"
)

// MasterPass defines struct for init
type MasterPass struct {
	db *gorm.DB
}

// New provides initialization of postgres for store master pass
func New(c *config.MasterPassBackend) (backend.MasterPassBackend, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Host, c.Port, c.User, c.DB, c.Password)
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	return &MasterPass{
		db: db,
	}, nil
}

// Get defines getting of master pass from backend
func (m *MasterPass) Get(key string) ([]byte, error) {
	return nil, nil
}

// Put defines putting of master pass to backend
func (m *MasterPass) Put(key string) error {
	return nil
}
