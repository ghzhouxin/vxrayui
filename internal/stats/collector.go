package stats

import (
	"zhouxin.learn/go/vxrayui/internal/logger"
	"zhouxin.learn/go/vxrayui/internal/storage"
	"zhouxin.learn/go/vxrayui/internal/types"
)

func GetSchemeYieldRate(scheme string) (*types.SchemeYieldRate, error) {
	rate, err := storage.Get[types.SchemeYieldRate](types.StorageKeySchemeYieldRate + scheme)
	if err != nil {
		logger.Error("failed to get scheme yield rate", "err", err.Error())
		return nil, err
	}
	return &rate, err
}

func SetSchemeYieldRate(rate *types.SchemeYieldRate) {
	err := storage.Set(types.StorageKeySchemeYieldRate, rate)
	if err != nil {
		logger.Error("failed to set scheme yield rate", "err", err.Error())
	}
}
