package main

import (
	"flag"

	"zhouxin.learn/go/vxrayui/config"
	"zhouxin.learn/go/vxrayui/internal/logger"
	"zhouxin.learn/go/vxrayui/internal/storage"
	"zhouxin.learn/go/vxrayui/internal/subscription"
)

func main() {
	flag.Parse()

	config.Init()
	logger.Init()
	storage.Init()

	sub := subscription.PickSubscription()
	logger.Info("Picked subscription", "scheme", sub.Scheme, "url", sub.Url)

	parser := subscription.NewSubscriptionParser()
	_ = parser.ParseSubscription(sub)

	/*
		engine := decision.NewEngine([]decision.Strategy{
			&decision.FreshnessStrategy{},
			&decision.SourcePriorityStrategy{},
		})

		stats := stats.NewCollector(store)
		sources := map[string]*subscription.SourceConfig{}

		// 启动轮询器
		poller := subscription.NewPoller(
			parser,
			store,
			engine,
			stats,
			sources,
		)
		poller.Run()
	*/
}
