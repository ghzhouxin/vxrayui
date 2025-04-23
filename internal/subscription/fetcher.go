package subscription

import (
	"bufio"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/xtls/libxray/share"
	"github.com/xtls/xray-core/infra/conf"
	"zhouxin.learn/go/vxrayui/config"
	"zhouxin.learn/go/vxrayui/internal/storage"
)

type Fetcher struct {
	client  *http.Client
	storage storage.Storage
}

func NewFetcher(client *http.Client, store storage.Storage) *Fetcher {
	return &Fetcher{client: client, storage: store}
}

func (f *Fetcher) Fetch(subscription config.Subsciption) []*conf.OutboundDetourConfig {
	resp, err := f.client.Get(subscription.URL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body
	if subscription.IsBase64 {
		body = base64.NewDecoder(base64.StdEncoding, resp.Body)
	}

	scanner := bufio.NewScanner(body)
	var outbounds []*conf.OutboundDetourConfig
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}

		link, err := url.Parse(text)
		if err != nil {
			slog.Error("parse url error", "url", text, "error", err)
			continue
		}

		shareLink := share.XrayShareLink{
			Link:    link,
			RawText: text,
		}

		if outbound, err := shareLink.Outbound(); err == nil {
			outbounds = append(outbounds, outbound)
		} else {
			slog.Error("shareLink.Outbound err", "url", text, "error", err)
			continue
		}
	}

	return outbounds
}

func (f *Fetcher) Validate(data []byte) bool {
	// 实现实际的验证逻辑
	return true
}
