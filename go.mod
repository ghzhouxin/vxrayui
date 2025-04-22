module zhouxin.learn/go/vxrayui

go 1.24.2

replace (
	github.com/xtls/libxray => ../libXray
	github.com/xtls/xray-core => ../Xray-core
)

require (
	go.etcd.io/bbolt v1.4.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
