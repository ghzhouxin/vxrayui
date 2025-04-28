package random

import (
	"math/rand"
	"sort"
)

func Pick[T any](items []T, weights []int) T {
	if len(items) <= 0 {
		panic("items cannot be empty")
	}

	if len(weights) < len(items) {
		weights = append(weights, make([]int, len(items)-len(weights))...)
	}

	prefixSums := make([]int, len(weights))
	for i, w := range weights {
		if w < 0 {
			w = 0 // 自动处理负权重
		}

		if i == 0 {
			prefixSums[i] = w
			continue
		}

		prefixSums[i] = prefixSums[i-1] + w
	}

	if prefixSums[len(prefixSums)-1] == 0 {
		return items[rand.Intn(len(items))]
	}

	return items[sort.SearchInts(prefixSums, rand.Intn(prefixSums[len(prefixSums)-1]))]
}
