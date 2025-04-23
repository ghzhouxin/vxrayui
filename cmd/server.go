package main

import (
	"zhouxin.learn/go/vxrayui/config"
	logger "zhouxin.learn/go/vxrayui/internal/log"
	"zhouxin.learn/go/vxrayui/internal/subscription"
)

func main() {
	config.Init()
	logger.Init()

	parser := subscription.NewSubscriptionParser()
	parser.ParseSubscription(config.Config.Subscriptions[0])

	/*
		store, err := storage.NewBoltStore("configs.db")
		if err != nil {
			log.Fatal(err)
		}

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
