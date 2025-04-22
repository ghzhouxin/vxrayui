package stats

import (
	"zhouxin.learn/go/vxrayui/internal/storage"
)

type Collector struct {
	storage storage.Storage
}

func NewCollector(store storage.Storage) *Collector {
	return &Collector{storage: store}
}

func (c *Collector) RecordValidation(configID string, isValid bool) error {
	// 实现统计记录逻辑
	return nil
}

func (c *Collector) GetValidityRate(configType string) float64 {
	// 实现统计查询逻辑
	return 0.95
}
