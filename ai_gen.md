# 主程序入口

cat > cmd/server.go <<'EOF'
package main

import (
"log"
"net/http"
"time"

    "zhouxin.learn/go/vxrayui/internal/config"
    "zhouxin.learn/go/vxrayui/internal/decision"
    "zhouxin.learn/go/vxrayui/internal/stats"
    "zhouxin.learn/go/vxrayui/internal/storage"

)

func main() {
store, err := storage.NewBoltStore("configs.db")
if err != nil {
log.Fatal(err)
}

    fetcher := config.NewFetcher(http.DefaultClient, store)
    engine := decision.NewEngine([]decision.Strategy{
    	&decision.FreshnessStrategy{},
    	&decision.SourcePriorityStrategy{},
    })
    statsCollector := stats.NewCollector(store)

    poller := config.NewPoller(
    	fetcher,
    	store,
    	engine,
    	statsCollector,
    )
    poller.Run(time.Minute)

}
EOF

# 存储层实现

cat > internal/storage/bolt_store.go <<'EOF'
package storage

import (
"encoding/json"
"time"

    "github.com/boltdb/bolt"

)

type BoltStore struct {
db \*bolt.DB
}

func NewBoltStore(path string) (_BoltStore, error) {
db, err := bolt.Open(path, 0600, &bolt.Options{
Timeout: 1 _ time.Second,
NoGrowSync: false,
FreelistType: bolt.FreelistArrayType,
})
if err != nil {
return nil, err
}

    err = db.Update(func(tx *bolt.Tx) error {
    	buckets := []string{"configs", "stats", "sources"}
    	for _, name := range buckets {
    		if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
    			return err
    		}
    	}
    	return nil
    })

    return &BoltStore{db: db}, err

}

func (s *BoltStore) StoreConfig(cfg *ConfigMetadata) error {
return s.db.Update(func(tx \*bolt.Tx) error {
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
err := s.db.View(func(tx \*bolt.Tx) error {
b := tx.Bucket([]byte("configs"))
data := b.Get([]byte(id))
return json.Unmarshal(data, &cfg)
})
return &cfg, err
}

type ConfigMetadata struct {
ID string
Content []byte
Version string
LastUpdated time.Time
Valid bool
SourceURL string
}
EOF

# 配置获取模块

cat > internal/config/fetcher.go <<'EOF'
package config

import (
"crypto/sha256"
"fmt"
"io"
"net/http"
"time"

    "zhouxin.learn/go/vxrayui/internal/storage"

)

type Fetcher struct {
client \*http.Client
storage storage.Storage
}

func NewFetcher(client *http.Client, store storage.Storage) *Fetcher {
return &Fetcher{
client: client,
storage: store,
}
}

func (f \*Fetcher) Fetch(url string) ([]byte, string, error) {
resp, err := f.client.Get(url)
if err != nil {
return nil, "", err
}
defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
    	return nil, "", err
    }

    hash := sha256.Sum256(data)
    return data, fmt.Sprintf("%x", hash), nil

}

func (f \*Fetcher) Validate(data []byte) bool {
// 实现实际的验证逻辑
return true
}
EOF

# 决策引擎

cat > internal/decision/engine.go <<'EOF'
package decision

import (
"zhouxin.learn/go/vxrayui/internal/storage"
)

type Engine struct {
strategies []Strategy
}

func NewEngine(strategies []Strategy) \*Engine {
return &Engine{strategies: strategies}
}

func (e *Engine) Decide(configs []*storage.ConfigMetadata) \*storage.ConfigMetadata {
scores := make(map[string]float64)

    for _, cfg := range configs {
    	for _, strat := range e.strategies  {
    		scores[cfg.ID] += strat.Score(cfg) * strat.Weight()
    	}
    }

    var best *storage.ConfigMetadata
    maxScore := -1.0
    for _, cfg := range configs {
    	if score := scores[cfg.ID]; score > maxScore {
    		best = cfg
    		maxScore = score
    	}
    }
    return best

}

type Strategy interface {
Score(cfg \*storage.ConfigMetadata) float64
Weight() float64
}
EOF

# 策略实现

cat > internal/decision/strategies.go <<'EOF'
package decision

import (
"time"

    "zhouxin.learn/go/vxrayui/internal/storage"

)

type FreshnessStrategy struct{}

func (s *FreshnessStrategy) Score(cfg *storage.ConfigMetadata) float64 {
age := time.Since(cfg.LastUpdated).Hours()
return 1 / (1 + age/24)
}

func (s \*FreshnessStrategy) Weight() float64 {
return 0.4
}

type SourcePriorityStrategy struct{}

func (s *SourcePriorityStrategy) Score(cfg *storage.ConfigMetadata) float64 {
// 根据 URL 判断源优先级
if strings.Contains(cfg.SourceURL, "prod") {
return 1.0
}
return 0.7
}

func (s \*SourcePriorityStrategy) Weight() float64 {
return 0.3
}
EOF

# 统计模块

cat > internal/stats/collector.go <<'EOF'
package stats

import (
"time"

    "zhouxin.learn/go/vxrayui/internal/storage"

)

type Collector struct {
storage storage.Storage
}

func NewCollector(store storage.Storage) \*Collector {
return &Collector{storage: store}
}

func (c \*Collector) RecordValidation(configID string, isValid bool) error {
// 实现统计记录逻辑
return nil
}

func (c \*Collector) GetValidityRate(configType string) float64 {
// 实现统计查询逻辑
return 0.95
}
EOF

# 实用工具

cat > pkg/utils/hashing.go <<'EOF'
package utils

import (
"crypto/sha256"
"fmt"
"io"
)

func CalculateHash(r io.Reader) (string, error) {
h := sha256.New()
if \_, err := io.Copy(h, r); err != nil {
return "", err
}
return fmt.Sprintf("%x", h.Sum(nil)), nil
}
EOF

# 创建 README

cat > README.md <<'EOF'

# Config Updater System

## Features

- Dynamic configuration management
- Smart polling algorithm
- Multi-strategy decision engine
- BoltDB-backed storage

## Quick Start

```sh
go run cmd/server.go
```

# Poller

cat > internal/subscription/poller.go <<'EOF'
package subscription

import (
"log"
"sync"
"time"

    "zhouxin.learn/go/vxrayui/internal/decision"
    "zhouxin.learn/go/vxrayui/internal/stats"
    "zhouxin.learn/go/vxrayui/internal/storage"

)

type Poller struct {
fetcher *Fetcher
storage storage.Storage
engine *decision.Engine
stats *stats.Collector
sources map[string]*SourceConfig
stopChan chan struct{}
wg sync.WaitGroup
}

type SourceConfig struct {
URL string
MinInterval time.Duration
MaxInterval time.Duration
LastCheck time.Time
FailureCount int
}

func NewPoller(
fetcher *Fetcher,
store storage.Storage,
engine *decision.Engine,
stats *stats.Collector,
sources map[string]*SourceConfig,
) \*Poller {
return &Poller{
fetcher: fetcher,
storage: store,
engine: engine,
stats: stats,
sources: sources,
stopChan: make(chan struct{}),
}
}

func (p \*Poller) Run() {
p.wg.Add(1)
defer p.wg.Done()

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
    	select {
    	case <-ticker.C:
    		p.pollAllSources()
    	case <-p.stopChan:
    		return
    	}
    }

}

func (p \*Poller) Stop() {
close(p.stopChan)
p.wg.Wait()
}

func (p \*Poller) pollAllSources() {
var wg sync.WaitGroup
semaphore := make(chan struct{}, 5) // 并发限制

    for url, source := range p.sources  {
    	if time.Since(source.LastCheck) < p.calculateInterval(source)  {
    		continue
    	}

    	wg.Add(1)
    	semaphore <- struct{}{}

    	go func(url string, source *SourceConfig) {
    		defer wg.Done()
    		defer func() { <-semaphore }()

    		p.pollSingleSource(url,  source)
    	}(url, source)
    }

    wg.Wait()

}

func (p *Poller) pollSingleSource(url string, source *SourceConfig) {
source.LastCheck = time.Now()

    data, hash, err := p.fetcher.Fetch(url)
    if err != nil {
    	source.FailureCount++
    	logger.Logger().Error("Failed to fetch %s (attempt %d): %v", url, source.FailureCount, err)
    	return
    }

    // 检查内容是否变化
    storedCfg, _ := p.storage.GetConfig(url)
    if storedCfg != nil && storedCfg.Valid {
    	if storedHash := sha256.Sum256(storedCfg.Content); fmt.Sprintf("%x", storedHash) == hash {
    		return // 内容未变化
    	}
    }

    // 验证并存储新配置
    if valid := p.fetcher.Validate(data);  valid {
    	newCfg := &storage.ConfigMetadata{
    		ID:         url,
    		Content:    data,
    		LastUpdated: time.Now(),
    		Valid:      true,
    		SourceURL:  url,
    	}

    	if err := p.storage.StoreConfig(newCfg);  err != nil {
    		logger.Logger().Error("Failed to store config from %s: %v", url, err)
    		return
    	}

    	source.FailureCount = 0 // 重置失败计数
    	p.stats.RecordValidation(url,  true)
    } else {
    	p.stats.RecordValidation(url,  false)
    }

}

func (p *Poller) calculateInterval(source *SourceConfig) time.Duration {
// 基于失败次数的指数退避
baseInterval := source.MinInterval
if source.FailureCount > 0 {
backoff := time.Duration(1 << min(source.FailureCount, 5)) \* time.Minute
baseInterval = min(source.MaxInterval, baseInterval+backoff)
}

    // 基于统计数据的动态调整
    validRate, _ := p.stats.GetValidityRate(source.URL)
    adjustment := 1.0 + (1.0 - validRate) // 有效性越低，检查越频繁
    return time.Duration(float64(baseInterval) * adjustment)

}

func min(a, b time.Duration) time.Duration {
if a < b {
return a
}
return b
}

# logging_files

cat > internal/logging/file_logger.go <<'EOF'
package logging

import (
"encoding/json"
"fmt"
"io"
"os"
"path/filepath"
"sync"
"time"
)

type FileLogger struct {
baseDir string
currentDay string
file \*os.File
mu sync.Mutex
}

func NewFileLogger(baseDir string) (\*FileLogger, error) {
if err := os.MkdirAll(baseDir, 0755); err != nil {
return nil, fmt.Errorf("failed to create log directory: %w", err)
}

    logger := &FileLogger{
    	baseDir: baseDir,
    }

    if err := logger.rotateIfNeeded();  err != nil {
    	return nil, err
    }

    return logger, nil

}

func (l \*FileLogger) rotateIfNeeded() error {
today := time.Now().Format("2006-01-02")

    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentDay  == today && l.file  != nil {
    	return nil
    }

    // 关闭旧文件
    if l.file  != nil {
    	_ = l.file.Close()
    }

    // 创建新文件
    logPath := filepath.Join(l.baseDir,  fmt.Sprintf("config_updater_%s.log",  today))
    f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
    	return fmt.Errorf("failed to open log file: %w", err)
    }

    l.currentDay  = today
    l.file  = f
    return nil

}

func (l \*FileLogger) WriteLog(level, message string, fields map[string]interface{}) error {
if err := l.rotateIfNeeded(); err != nil {
return err
}

    logEntry := map[string]interface{}{
    	"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
    	"level":     level,
    	"message":   message,
    	"fields":    fields,
    }

    entryBytes, err := json.Marshal(logEntry)
    if err != nil {
    	return fmt.Errorf("failed to marshal log entry: %w", err)
    }

    l.mu.Lock()
    defer l.mu.Unlock()

    if _, err := l.file.Write(append(entryBytes,  '\n')); err != nil {
    	return fmt.Errorf("failed to write log: %w", err)
    }

    return nil

}

func (l _FileLogger) Cleanup(maxAge time.Duration) error {
files, err := filepath.Glob(filepath.Join(l.baseDir, "config*updater*_.log"))
if err != nil {
return err
}

    now := time.Now()
    for _, file := range files {
    	info, err := os.Stat(file)
    	if err != nil {
    		continue
    	}

    	if now.Sub(info.ModTime()) > maxAge {
    		if err := os.Remove(file); err != nil {
    			return err
    		}
    	}
    }
    return nil

}

func (l \*FileLogger) Close() error {
l.mu.Lock()
defer l.mu.Unlock()

    if l.file  != nil {
    	return l.file.Close()
    }
    return nil

}

# log

cat > pkg/logger/logger.go <<'EOF'
package logger

import (
"context"
"time"

    "github.com/yourproject/config-updater/internal/logging"

)

type Logger struct {
fileLogger \*logging.FileLogger
service string
}

func NewLogger(baseDir, service string) (\*Logger, error) {
fl, err := logging.NewFileLogger(baseDir)
if err != nil {
return nil, err
}

    return &Logger{
    	fileLogger: fl,
    	service:    service,
    }, nil

}

func (l \*Logger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
l.log(ctx, "INFO", msg, fields)
}

func (l \*Logger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
l.log(ctx, "ERROR", msg, fields)
}

func (l \*Logger) log(ctx context.Context, level, msg string, fields map[string]interface{}) {
if fields == nil {
fields = make(map[string]interface{})
}

    fields["service"] = l.service
    if traceID := ctx.Value("trace_id"); traceID != nil {
    	fields["trace_id"] = traceID
    }

    _ = l.fileLogger.WriteLog(level,  msg, fields) // 忽略错误避免循环

}

func (l \*Logger) StartCleanup(interval, maxAge time.Duration) {
go func() {
ticker := time.NewTicker(interval)
defer ticker.Stop()

    	for range ticker.C {
    		_ = l.fileLogger.Cleanup(maxAge)
    	}
    }()

}

func (l \*Logger) Close() error {
return l.fileLogger.Close()
}

# main

cat > cmd/server.go <<'EOF'
package main

import (
"context"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

    "github.com/yourproject/config-updater/internal/config"
    "github.com/yourproject/config-updater/internal/decision"
    "github.com/yourproject/config-updater/internal/stats"
    "github.com/yourproject/config-updater/internal/storage"
    "github.com/yourproject/config-updater/pkg/logger"

)

func main() {
// 初始化日志系统
appLogger, err := logger.NewLogger("/var/log/config-updater", "config-updater")
if err != nil {
log.Fatalf("Failed to init logger: %v", err)
}
defer appLogger.Close()
appLogger.StartCleanup(24*time.Hour, 30*24\*time.Hour) // 每天清理，保留 30 天

    ctx := context.WithValue(context.Background(), "trace_id", time.Now().UnixNano())

    // 初始化其他组件
    store, err := storage.NewBoltStore("configs.db")
    if err != nil {
    	appLogger.Error(ctx, "Failed to init storage", map[string]interface{}{
    		"error": err.Error(),
    	})
    	os.Exit(1)
    }

    // 其他初始化代码...
    appLogger.Info(ctx, "Service starting up", nil)

    // 优雅退出处理
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    appLogger.Info(ctx, "Shutting down gracefully", nil)
    time.Sleep(1 * time.Second) // 确保日志写完

}
