package decision

import (
	"zhouxin.learn/go/vxrayui/internal/storage"
)

type Engine struct {
	strategies []Strategy
}

type Strategy interface {
	Score(cfg *storage.ConfigMetadata) float64
	Weight() float64
}

func NewEngine(strategies []Strategy) *Engine {
	return &Engine{strategies: strategies}
}

func (e *Engine) Decide(configs []*storage.ConfigMetadata) *storage.ConfigMetadata {
	scores := make(map[string]float64)
	for _, cfg := range configs {
		for _, strat := range e.strategies {
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
