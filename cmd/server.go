package main

import (
	"log"
	"net/http"

	"zhouxin.learn/go/vxrayui/internal/decision"
	"zhouxin.learn/go/vxrayui/internal/stats"
	"zhouxin.learn/go/vxrayui/internal/storage"
	"zhouxin.learn/go/vxrayui/internal/subscription"
)

func main() {
	// 初始化存储
	store, err := storage.NewBoltStore("configs.db")
	if err != nil {
		log.Fatal(err)
	}

	// 初始化各模块
	fetcher := subscription.NewFetcher(http.DefaultClient, store)
	engine := decision.NewEngine([]decision.Strategy{
		&decision.FreshnessStrategy{},
		&decision.SourcePriorityStrategy{},
	})
	stats := stats.NewCollector(store)

	sources := map[string]*subscription.SourceConfig{}

	// 启动轮询器
	poller := subscription.NewPoller(
		fetcher,
		store,
		engine,
		stats,
		sources,
	)
	poller.Run()
}
