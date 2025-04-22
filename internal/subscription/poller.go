package subscription

import (
	"crypto/sha256"
	"fmt"
	"log"
	"sync"
	"time"

	"zhouxin.learn/go/vxrayui/internal/decision"
	"zhouxin.learn/go/vxrayui/internal/stats"
	"zhouxin.learn/go/vxrayui/internal/storage"
)

type Poller struct {
	fetcher  *Fetcher
	storage  storage.Storage
	engine   *decision.Engine
	stats    *stats.Collector
	sources  map[string]*SourceConfig
	stopChan chan struct{}
	wg       sync.WaitGroup
}

type SourceConfig struct {
	URL          string
	MinInterval  time.Duration
	MaxInterval  time.Duration
	LastCheck    time.Time
	FailureCount int
}

func NewPoller(
	fetcher *Fetcher,
	store storage.Storage,
	engine *decision.Engine,
	stats *stats.Collector,
	sources map[string]*SourceConfig,
) *Poller {
	return &Poller{
		fetcher:  fetcher,
		storage:  store,
		engine:   engine,
		stats:    stats,
		sources:  sources,
		stopChan: make(chan struct{}),
	}
}

func (p *Poller) Run() {
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

func (p *Poller) Stop() {
	close(p.stopChan)
	p.wg.Wait()
}

func (p *Poller) pollAllSources() {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // 并发限制

	for url, source := range p.sources {
		if time.Since(source.LastCheck) < p.calculateInterval(source) {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(url string, source *SourceConfig) {
			defer wg.Done()
			defer func() { <-semaphore }()

			p.pollSingleSource(url, source)
		}(url, source)
	}

	wg.Wait()
}

func (p *Poller) pollSingleSource(url string, source *SourceConfig) {
	source.LastCheck = time.Now()

	data, hash, err := p.fetcher.Fetch(url)
	if err != nil {
		source.FailureCount++
		log.Printf("Failed to fetch %s (attempt %d): %v", url, source.FailureCount, err)
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
	if valid := p.fetcher.Validate(data); valid {
		newCfg := &storage.ConfigMetadata{
			ID:          url,
			Content:     data,
			LastUpdated: time.Now(),
			Valid:       true,
			SourceURL:   url,
		}

		if err := p.storage.StoreConfig(newCfg); err != nil {
			log.Printf("Failed to store config from %s: %v", url, err)
			return
		}

		source.FailureCount = 0 // 重置失败计数
		p.stats.RecordValidation(url, true)
	} else {
		p.stats.RecordValidation(url, false)
	}
}

func (p *Poller) calculateInterval(source *SourceConfig) time.Duration {
	// 基于失败次数的指数退避
	baseInterval := source.MinInterval
	if source.FailureCount > 0 {
		backoff := time.Duration(1<<min(time.Duration(source.FailureCount), 5)) * time.Minute
		baseInterval = min(source.MaxInterval, baseInterval+backoff)
	}

	// 基于统计数据的动态调整
	validRate := p.stats.GetValidityRate(source.URL)
	adjustment := 1.0 + (1.0 - validRate) // 有效性越低，检查越频繁
	return time.Duration(float64(baseInterval) * adjustment)
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
