package subscription

import (
	"zhouxin.learn/go/vxrayui/config"
	"zhouxin.learn/go/vxrayui/internal/stats"
	"zhouxin.learn/go/vxrayui/pkg/random"
)

func PickSubscription() *config.Subscription {
	var subs []*config.Subscription
	var weights []int
	for _, sub := range config.GetSubscriptions() {
		if !sub.Enabled {
			continue
		}

		scheme := sub.Scheme
		if scheme == "" {
			scheme = config.DefalutScheme
		}

		yieldRate, _ := stats.GetSchemeYieldRate(sub.Scheme)
		subs = append(subs, sub)
		if yieldRate == nil || yieldRate.Total == 0 {
			weights = append(weights, 0)
		} else {
			weights = append(weights, yieldRate.Yiled*10000/yieldRate.Total)
		}
	}

	return random.Pick(subs, weights)
}
