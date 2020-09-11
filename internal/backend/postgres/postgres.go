package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/saromanov/mystery/config"
	"github.com/saromanov/mystery/internal/backend"
	"github.com/saromanov/mystery/internal/crypto"
)

// Postgres defines backend for postgres
type Postgres struct {
	db *gorm.DB
}

// Mystery defines structure for store in Postgres
type Mystery struct {
	ID             uint64    `gorm:"primaryKey;AUTO_INCREMENT;NOT NULL"`
	Key            string    `gorm:"NOT NULL"`
	Value          []byte    `gorm:"NOT NULL"`
	UserID         string    `gorm:"index"`
	CreatedAt      time.Time `gorm:"NOT NULL"`
	CurrentVersion uint      `gorm:"NOT NULL;default:0"`
	MaxVersion     uint      `gorm:"NOT NULL;default:0"`
	ExpiredAfter   *time.Duration
	UpdatedAt      time.Time
}

// New provides initialization of postgres for store master pass
func New(c *config.Config) (backend.Backend, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Backend.Host, c.Backend.Port, c.Backend.User, c.Backend.DB, c.Backend.Password)
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	db.AutoMigrate(&Mystery{})
	return &Postgres{
		db: db,
	}, nil
}

// Get defines getting a secret from backend
func (m *Postgres) Get(masterKey, key []byte) (backend.Secret, error) {
	r, err := m.get(masterKey, key)
	if err != nil {
		return backend.Secret{}, fmt.Errorf("unable to get secret: %v", err)
	}
	decrypted, err := crypto.DecryptAES(masterKey, r.Value)
	if err != nil {
		return backend.Secret{}, fmt.Errorf("get: unable to decrypt value: %v", err)
	}
	return backend.Secret{
		Key:   key,
		Value: decrypted,
	}, nil
}

func (m *Postgres) get(masterKey, key []byte) (Mystery, error) {
	var r Mystery
	if err := m.db.Find(&r, &Mystery{
		Key: string(key),
	}).Error; err != nil {
		return r, fmt.Errorf("unable to get secret: %v", err)
	}
	if expired := checkExpired(r); expired {
		return r, fmt.Errorf("data with key %s has expired", string(key))
	}
	return r, nil
}

// checkExpired provides checking is value has expired
func checkExpired(r Mystery) bool {
	if r.ExpiredAfter == nil {
		return false
	}
	now := time.Now().UTC()
	if now.Sub(r.CreatedAt) > *r.ExpiredAfter {
		return true
	}
	return false

}

// Put defines putting a secret to backend
func (m *Postgres) Put(masterKey []byte, secret backend.Secret) error {
	encryptedValue, err := crypto.EncryptAES(masterKey, secret.Value)
	if err != nil {
		return fmt.Errorf("put: unable to encrypt data: %v", err)
	}
	m.db.Create(&Mystery{
		Key:          string(secret.Key),
		Value:        encryptedValue,
		CreatedAt:    time.Now().UTC(),
		ExpiredAfter: secret.ExpiredAfter,
	})
	return nil
}

func (m *Postgres) Update(masterKey []byte, secret backend.Secret) error {
	data, err := m.get(masterKey, secret.Key)
	if err != nil {
		return err
	}
	data.CurrentVersion++
	return m.db.Update(data)
}
