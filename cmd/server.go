package main

import (
	"flag"

	"zhouxin.learn/go/vxrayui/config"
	"zhouxin.learn/go/vxrayui/internal/logger"
	"zhouxin.learn/go/vxrayui/internal/subscription"
)

func main() {
	flag.Parse()

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
