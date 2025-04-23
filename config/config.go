package config

import (
	"log"
	"os"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	Level   string `yaml:"level" json:"level"`
	Console struct {
		Enabled bool   `yaml:"enabled" json:"enabled"`
		Format  string `yaml:"format" json:"format"`
	} `yaml:"console" json:"console"`
	File struct {
		Enabled    bool   `yaml:"enabled" json:"enabled"`
		Path       string `yaml:"path" json:"path"`
		Filename   string `yaml:"filename" json:"filename"`
		MaxSize    int    `yaml:"max_size" json:"max_size"`
		MaxAge     int    `yaml:"max_age" json:"max_age"`
		MaxBackups int    `yaml:"max_backups" json:"max_backups"`
		ShardBy    string `yaml:"shard_by" json:"shard_by"`
		Compress   bool   `yaml:"compress" json:"compress"`
	} `yaml:"logger" json:"logger"`
}

type Subsciption struct {
	Name     string `json:"name" yaml:"name"`
	URL      string `json:"url" yaml:"url"`
	IsBase64 bool   `json:"is_base64" yaml:"is_base64"`
}

type config struct {
	Log          *LogConfig    `json:"logger" yaml:"logger"`
	Subsciptions []Subsciption `json:"subsciptions" yaml:"subsciptions"`
}

var globalConfig atomic.Value

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	globalConfig.Store(&cfg)
}

func GetConfig() *config {
	return globalConfig.Load().(*config)
}
