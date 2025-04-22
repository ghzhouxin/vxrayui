package storage

import (
	"encoding/json"
	"time"

	"go.etcd.io/bbolt"
)

type BoltStore struct {
	db *bbolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("configs"))
		return err
	})
	return &BoltStore{db: db}, err
}

func (s *BoltStore) StoreConfig(cfg *ConfigMetadata) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("configs"))
		data, err := json.Marshal(cfg)
		if err != nil {
			return err
		}
		return b.Put([]byte(cfg.ID), data)
	})
}

func (s *BoltStore) GetConfig(id string) (*ConfigMetadata, error) {
	var cfg ConfigMetadata
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("configs"))
		data := b.Get([]byte(id))
		return json.Unmarshal(data, &cfg)
	})
	return &cfg, err
}
