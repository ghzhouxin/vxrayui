package decision

import (
	"strings"
	"time"

	"zhouxin.learn/go/vxrayui/internal/storage"
)

type FreshnessStrategy struct{}

func (s *FreshnessStrategy) Score(cfg *storage.ConfigMetadata) float64 {
	age := time.Since(cfg.LastUpdated).Hours()
	return 1 / (1 + age/24)
}

func (s *FreshnessStrategy) Weight() float64 {
	return 0.4
}

type SourcePriorityStrategy struct{}

func (s *SourcePriorityStrategy) Score(cfg *storage.ConfigMetadata) float64 {
	// 根据URL判断源优先级
	if strings.Contains(cfg.SourceURL, "prod") {
		return 1.0
	}
	return 0.7
}

func (s *SourcePriorityStrategy) Weight() float64 {
	return 0.3
}
