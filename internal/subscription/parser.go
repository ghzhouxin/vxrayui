package subscription

import (
	"bufio"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xtls/xray-core/infra/conf"
	"zhouxin.learn/go/vxrayui/config"
	"zhouxin.learn/go/vxrayui/internal/logger"
	"zhouxin.learn/go/vxrayui/internal/types"
	"zhouxin.learn/go/vxrayui/pkg/counter"
	"zhouxin.learn/go/vxrayui/pkg/xray"
)

// SubscriptionParser 用于解析订阅内容并生成 OutboundDetourConfig
type SubscriptionParser struct{}

// NewSubscriptionParser 创建一个新的 SubscriptionParser
func NewSubscriptionParser() *SubscriptionParser {
	return &SubscriptionParser{}
}

// ParseSubscription 从订阅链接拉取内容并解析为 OutboundDetourConfig
func (p *SubscriptionParser) ParseSubscription(subscription *config.Subscription) []*conf.OutboundDetourConfig {
	resp, err := http.Get(subscription.Url)
	if err != nil {
		logger.Error("Failed to fetch subscription", "err", err.Error())
		return nil
	}
	defer resp.Body.Close()

	reader := decodeBody(resp.Body, subscription.IsBase64)
	if reader == nil {
		logger.Error("Failed to decode subscription body")
		return nil
	}

	return parseSubscriptionContent(reader)
}

// decodeBody 根据是否 base64 解码返回正确的 reader
func decodeBody(body io.Reader, isBase64 bool) io.Reader {
	if isBase64 {
		return base64.NewDecoder(base64.StdEncoding, body)
	}
	return body
}

// parseSubscriptionContent 解析订阅内容
func parseSubscriptionContent(reader io.Reader) []*conf.OutboundDetourConfig {
	scanner := bufio.NewScanner(reader)
	var outbounds []*conf.OutboundDetourConfig

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if !isValidLink(line) {
			logger.Error("Unsupported subscription", "url", line)
			counter.Incr("subscription.invalid", 1)
			continue
		}

		// TODO 去重

		// TODO 抽象方法
		link, err := url.Parse(line)
		if err != nil {
			logger.Error("Invalid Url in subscription", "url", line, "err", err.Error())
			counter.Incr("subscription.invalid", 1)
			continue
		}
		shareLink := xray.XrayShareLink{
			Link:    link,
			RawText: line,
		}
		outbound, err := shareLink.Outbound()
		if err != nil {
			logger.Error("Failed to parse outbound from link", "link", line, "err", err.Error())
			counter.Incr("subscription.parse.error", 1)
			continue
		}

		// string(*outbound.StreamSetting.Network) + outbound.StreamSetting.Security
		outbounds = append(outbounds, outbound)
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error reading subscription", "err", err.Error())
	}

	cnt := counter.Get("subscription.invalid") + counter.Get("subscription.parse.error")
	logger.Info("Parsed outbounds from subscription result", "total", len(outbounds), "invalid", cnt)
	return outbounds
}

func (p *SubscriptionParser) Fetch(url string) ([]byte, string, error) {
	return nil, "", nil
}

func (p *SubscriptionParser) Validate(data []byte) bool {
	return true
}

// isValidLink 判断是否为支持的协议链接
func isValidLink(link string) bool {
	for _, scheme := range types.SupportedSchemes {
		if strings.HasPrefix(link, scheme.String()) {
			return true
		}
	}
	return false
}
