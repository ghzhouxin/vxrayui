package config

import (
	"embed"
	"flag"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Logger struct {
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
	} `yaml:"file" json:"file"`
}

type Subscription struct {
	Name     string `json:"name" yaml:"name"`
	Url      string `json:"url" yaml:"url"`
	IsBase64 bool   `json:"is_base64" yaml:"is_base64"`
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Scheme   string `json:"scheme" yaml:"scheme"`
}

type Storage struct {
	Type string `json:"type" yaml:"type"`
	Path string `json:"path" yaml:"path"`
}

type config struct {
	Logger        *Logger         `json:"logger" yaml:"logger"`
	Subscriptions []*Subscription `json:"subscriptions" yaml:"subscriptions"`
	Storage       *Storage        `json:"storage" yaml:"storage"`
}

const DefalutScheme string = "mix"

var (
	//go:embed config.yaml
	defaultConfigFile embed.FS
	configFilePath    string = *flag.String("config", "", "config file path")
	initOnce          sync.Once

	cfg *config
)

func GetLogger() *Logger {
	return cfg.Logger
}

func GetSubscriptions() []*Subscription {
	return cfg.Subscriptions
}

func GetStorage() *Storage {
	return cfg.Storage
}

func Init() {
	initOnce.Do(func() {
		initConfig()
	})
}

func initConfig() {
	var configData []byte
	if configFilePath != "" {
		if data, err := os.ReadFile(configFilePath); err == nil {
			configData = data
		} else {
			log.Fatalf("failed to read config file: %v", err)
		}
	} else {
		if data, err := defaultConfigFile.ReadFile("config.yaml"); err == nil {
			configData = data
		} else {
			log.Fatalf("failed to read default config file: %v", err)
		}
	}

	var config config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	cfg = &config
}
