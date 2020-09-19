package postgres

import (
	"errors"
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
	Namespace      string    `gorm:"NOT NULL"`
	Data           []byte    `gorm:"NOT NULL"`
	UserID         string    `gorm:"index"`
	CreatedAt      time.Time `gorm:"NOT NULL"`
	CurrentVersion uint64    `gorm:"NOT NULL;default:0"`
	MaxVersion     uint64    `gorm:"NOT NULL;default:0"`
	Compressed     bool
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
func (m *Postgres) Get(masterKey, namespace []byte) (backend.Secret, error) {
	r, err := m.get(namespace)
	if err != nil {
		return backend.Secret{}, fmt.Errorf("unable to get secret: %v", err)
	}
	decrypted, err := crypto.DecryptAES(masterKey, r.Data)
	if err != nil {
		return backend.Secret{}, fmt.Errorf("get: unable to decrypt value: %v", err)
	}
	return backend.Secret{
		Namespace:  namespace,
		Data:       decrypted,
		Compressed: r.Compressed,
	}, nil
}

// Delete defines deleting secret from backend
func (m *Postgres) Delete(masterKey, namespace []byte) error {
	if len(namespace) == 0 {
		return errors.New("delete: namespace is not defined")
	}
	if err := m.db.Where("namespace = ?", string(namespace)).Delete(&Mystery{}).Error; err != nil {
		return fmt.Errorf("unable to delete secret: %v", err)
	}
	return nil
}

// get data by the key
func (m *Postgres) get(key []byte) (Mystery, error) {
	var r Mystery
	count, err := m.countByKey(string(key))
	if err != nil {
		return r, err
	}
	if err := m.db.Find(&r, &Mystery{
		Namespace:      string(key),
		CurrentVersion: count,
	}).Error; err != nil {
		return r, fmt.Errorf("unable to get secret: %v", err)
	}
	if expired := checkExpired(r); expired {
		return r, fmt.Errorf("data with key %s has expired", string(key))
	}
	return r, nil
}

// countByKey returns number of secrets by the key
func (m *Postgres) countByKey(key string) (uint64, error) {
	var count uint64
	if err := m.db.Model(&Mystery{}).Where("namespace = ?", key).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("unable to get count of keys: %v", err)
	}
	return count, nil
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
	encryptedValue, err := crypto.EncryptAES(masterKey, secret.Data)
	if err != nil {
		return fmt.Errorf("put: unable to encrypt data: %v", err)
	}
	err = m.db.Create(&Mystery{
		Namespace:      string(secret.Namespace),
		Data:           encryptedValue,
		CreatedAt:      time.Now().UTC(),
		ExpiredAfter:   secret.ExpiredAfter,
		CurrentVersion: 1,
		MaxVersion:     1,
		Compressed:     secret.Compressed,
	}).Error
	if err != nil {
		return fmt.Errorf("unable to create data: %v", err)
	}
	return nil
}

// Update provides updating of data
func (m *Postgres) Update(masterKey []byte, secret backend.Secret) error {
	data, err := m.get(secret.Namespace)
	if err != nil {
		return err
	}
	data.CurrentVersion++
	data.Data = secret.Data
	data.UpdatedAt = time.Now().UTC()
	err = m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(Mystery{}).Update(&Mystery{
			Data:       secret.Data,
			Compressed: secret.Compressed,
			MaxVersion: data.MaxVersion + 1,
		}).Error; err != nil {
			return fmt.Errorf("unable to update old value: %v", err)
		}

		encryptedValue, err := crypto.EncryptAES(masterKey, secret.Data)
		if err != nil {
			return fmt.Errorf("put: unable to encrypt data: %v", err)
		}
		if err := m.db.Create(&Mystery{
			Namespace:      string(secret.Namespace),
			Data:           encryptedValue,
			CreatedAt:      time.Now().UTC(),
			ExpiredAfter:   secret.ExpiredAfter,
			CurrentVersion: data.MaxVersion + 1,
			MaxVersion:     data.MaxVersion + 1,
			Compressed:     secret.Compressed,
		}).Error; err != nil {
			return fmt.Errorf("unable to create new record: %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
