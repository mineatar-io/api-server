package store

import (
	"image"
	"time"
)

var (
	StoreTypes map[string]Store = map[string]Store{
		"filestore": &FileStore{},
	}
)

type Store interface {
	Initialize(config map[string]interface{}) error
	GetBytes(id string) ([]byte, bool, error)
	GetNRGBA(id string) (*image.NRGBA, bool, error)
	Exists(id string) (bool, error)
	SetBytes(id string, data []byte, ttl time.Duration) error
	Delete(id string) error
	Close() error
}
