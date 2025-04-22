package storage

import "time"

type ConfigMetadata struct {
	ID          string
	Content     []byte
	Version     string
	LastUpdated time.Time
	Valid       bool
	SourceURL   string
}

type Storage interface {
	StoreConfig(cfg *ConfigMetadata) error
	GetConfig(id string) (*ConfigMetadata, error)
}
