package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"time"

	vxnet "github.com/xtls/xray-core/common/net"
	vxcore "github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
)

type MeasureResult int

const (
	None MeasureResult = iota
	Select
	Delete
)

const (
	curlUrl = "http://www.gstatic.com/generate_204"
)

func OutboundMeasure(outbound *conf.OutboundDetourConfig) (MeasureResult, error) {
	if outbound == nil {
		return None, nil
	}

	vxrayConfig := conf.Config{
		OutboundConfigs: []conf.OutboundDetourConfig{*outbound},
	}
	vxrayConfigPb, err := vxrayConfig.Build()
	if err != nil {
		return None, err
	}

	vxrayInstance, err := vxcore.New(vxrayConfigPb)
	if err := vxrayInstance.Start(); err != nil {
		return None, err
	}
	defer vxrayInstance.Close()

	client := BuildProxyClient(vxrayInstance)

	times := 3
	res := []int64{}
	for i := 0; i < times; i++ {
		msec, err := OutBoundCurlDelay(client)
		if err != nil && !os.IsTimeout(err) {
			fmt.Println(err)
			break
		}
		fmt.Println("delay: ", msec)
		res = append(res, msec)
	}

	if len(res) < times {
		return None, err
	}
	d, s := 0, 0
	for i := 0; i < times; i++ {
		if res[i] > 3000 {
			d++
		}
		if res[i] <= 1000 {
			s++
		}
	}
	if d >= times {
		return Delete, err
	}
	if s >= times {
		return Select, err
	}
	return None, err
}

func BuildProxyClient(vxrayInstance *vxcore.Instance) *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 3 * time.Second,
			DisableKeepAlives:   true,

			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dest, err := vxnet.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
				if err != nil {
					return nil, err
				}
				return vxcore.Dial(ctx, vxrayInstance, dest)
			},
		},
	}
}

func OutBoundCurlDelay(client *http.Client) (delayMillisecond int64, err error) {
	startTime := time.Now()
	resp, err := client.Get(curlUrl)
	if err != nil {
		fmt.Println("http err: ", err)
		return math.MaxInt64, err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io err: ", err)
		return math.MaxInt64, err
	}

	return time.Since(startTime).Milliseconds(), nil
}
