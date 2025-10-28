package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type VKVideoFetcher struct {
	URL string
}

type MovieDataFetcher interface {
	FetchMoviesData(ctx context.Context) ([]byte, error)
}

func NewVKVideoFetcher(url string) *VKVideoFetcher {
	return &VKVideoFetcher{URL: strings.TrimSpace(url)}
}

func (f *VKVideoFetcher) FetchMoviesData(ctx context.Context) ([]byte, error) {
	body := bytes.NewBufferString("al=1&silent_loading=1")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	utf8Data, err := charmap.Windows1251.NewDecoder().Bytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response from windows-1251 to UTF-8: %w", err)
	}

	return utf8Data, nil
}
