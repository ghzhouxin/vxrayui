package config

import (
	_ "embed"
	"log"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configData []byte

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
	} `yaml:"file" json:"file"`
}

type Subscription struct {
	Name     string `json:"name" yaml:"name"`
	Url      string `json:"url" yaml:"url"`
	IsBase64 bool   `json:"is_base64" yaml:"is_base64"`
}

type config struct {
	Logger        *LogConfig      `json:"logger" yaml:"logger"`
	Subscriptions []*Subscription `json:"subscriptions" yaml:"subscriptions"`
}

var (
	Config   *config
	initOnce sync.Once
)

func Init() {
	initOnce.Do(func() {
		initConfig()
	})
}

func initConfig() {
	// TODO 热加载外部文件
	// if data, err := os.ReadFile("custom.yaml");  err == nil {
	// 	// 使用外部配置
	// } else {
	// 	// 2. 回退到嵌入配置
	// 	data, _ = configFile.ReadFile("default.yaml")
	// }

	// data, err := configFile.ReadFile("config.yaml")
	// if err != nil {
	// 	log.Fatalf("failed to read config file: %v", err)
	// }

	var cfg config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	Config = &cfg
}
