package subscription

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"zhouxin.learn/go/vxrayui/internal/storage"
)

type Fetcher struct {
	client  *http.Client
	storage storage.Storage
}

func NewFetcher(client *http.Client, store storage.Storage) *Fetcher {
	return &Fetcher{client: client, storage: store}
}

func (f *Fetcher) Fetch(url string) ([]byte, string, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	hash := sha256.Sum256(data)
	return data, fmt.Sprintf("%x", hash), nil
}

func (f *Fetcher) Validate(data []byte) bool {
	// 实现实际的验证逻辑
	return true
}
