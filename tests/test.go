package main

import (
	"fmt"
	"net/url"

	"zhouxin.learn/go/vxrayui/pkg/xray"
)

// var ss string = "ss://eyJhZGQiOiJodHRwczovL2dpdGh1Yi5jb20vQUxJSUxBUFJPL3YycmF5TkctQ29uZmlnIiwiYWlkIjoiMCIsImFscG4iOiIiLCJmcCI6IiIsImhvc3QiOiIiLCJpZCI6IkZyZWUiLCJuZXQiOiJ0Y3AiLCJwYXRoIjoiIiwicG9ydCI6IjQzMyIsInBzIjoi8J+SgPCfmI4gUHJvamVjdCBCeSBBTElJTEFQUk8iLCJzY3kiOiJjaGFjaGEyMC1wb2x5MTMwNSIsInNuaSI6IiIsInRscyI6IiIsInR5cGUiOiJub25lIiwidiI6IjIifQ=="
var ss string = "ss://8dc5b94a-382f-4d34-b44e-d78ba12aee1a@www.speedtest.net:8880?path=%2FChannel----VPNCUSTOMIZE----VPNCUSTOMIZE----VPNCUSTOMIZE---VPNCUSTOMIZE---VPNCUSTOMIZE---VPNCUSTOMIZE%3Fed%3D2048&security=none&alpn=h3%2Ch2%2Chttp%2F1.1&encryption=none&host=join.VPNCUSTOMIZE.iran.ir.arvancloud.ir.nett.ddns-ip.net.&fp=randomized&type=ws#👉🆔@v2ray_configs_pool📡🇨🇦®️Canada©️Toronto🅿️ping:15.62ms`"

func main() {
	link, _ := url.Parse(ss)

	shareLink := xray.XrayShareLink{
		Link:    link,
		RawText: ss,
	}
	outbound, _ := shareLink.Outbound()
	fmt.Println("outbound", outbound)
}
