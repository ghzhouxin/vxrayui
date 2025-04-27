package subscription

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xtls/libxray/share"
	"github.com/xtls/xray-core/infra/conf"
	"zhouxin.learn/go/vxrayui/config"
	"zhouxin.learn/go/vxrayui/internal/logger"
	"zhouxin.learn/go/vxrayui/pkg/counter"
)

// 支持的协议前缀
var supportedSchemes = []string{
	"vless://",
	"vmess://",
	"socks://",
	"ss://",
	"trojan://",
}

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
		logger.Logger.Error(fmt.Sprintf("Failed to fetch subscription: %v", err))
		return nil
	}
	defer resp.Body.Close()

	reader := decodeBody(resp.Body, subscription.IsBase64)
	if reader == nil {
		logger.Logger.Error("Failed to decode subscription body")
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
			logger.Logger.Error(fmt.Sprintf("Unsupported subscription: %s", line))
			counter.Incr("subscription.invalid", 1)
			continue
		}

		// TODO 去重

		// TODO 抽象方法
		link, err := url.Parse(line)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Invalid URL in subscription: %s, error: %v", line, err))
			counter.Incr("subscription.invalid", 1)
			continue
		}
		shareLink := share.XrayShareLink{
			Link:    link,
			RawText: line,
		}
		outbound, err := shareLink.Outbound()
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to parse outbound from link: %s, error: %v", line, err))
			counter.Incr("subscription.parse.error", 1)
			continue
		}
		outbounds = append(outbounds, outbound)
	}

	if err := scanner.Err(); err != nil {
		logger.Logger.Error(fmt.Sprintf("Error reading subscription: %v", err))
	}

	cnt := counter.Get("subscription.invalid") + counter.Get("subscription.parse.error")
	logger.Logger.Info(fmt.Sprintf("Parsed %d outbounds from subscription, error cnt: %v", len(outbounds), cnt))
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
	for _, scheme := range supportedSchemes {
		if strings.HasPrefix(link, scheme) {
			return true
		}
	}
	return false
}
