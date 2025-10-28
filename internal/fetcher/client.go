package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type VKVideoClient struct {
	BaseURL string
	client  *http.Client
}

func NewVKVideoClient(baseURL string) *VKVideoClient {
	return &VKVideoClient{
		BaseURL: strings.TrimSpace(baseURL),
		client:  &http.Client{},
	}
}

func (c *VKVideoClient) fetchNextPage(ctx context.Context, nextFrom string, showcaseSessionID int) ([]byte, error) {
	var bodyBuffer bytes.Buffer
	bodyBuffer.WriteString("al=1&silent_loading=1")
	if nextFrom != "" {
		bodyBuffer.WriteString("&next_from=")
		bodyBuffer.WriteString(nextFrom)
	}
	if showcaseSessionID != 0 {
		bodyBuffer.WriteString("&__showcase_session_id=")
		bodyBuffer.WriteString(strconv.Itoa(showcaseSessionID))
	}

	body := bytes.NewReader(bodyBuffer.Bytes())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, body)
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

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
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

func (c *VKVideoClient) extractInnerPayloadData(utf8Data []byte) (map[string]interface{}, error) {
	var response struct {
		Payload []interface{} `json:"payload"`
	}

	if err := json.Unmarshal(utf8Data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal root JSON for pagination: %w", err)
	}

	if len(response.Payload) < 2 {
		return nil, fmt.Errorf("invalid payload structure for pagination: expected at least 2 elements")
	}

	var innerPayloadData map[string]interface{}

	switch v := response.Payload[1].(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &innerPayloadData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inner JSON string: %w", err)
		}
	case map[string]interface{}:
		innerPayloadData = v
	case []interface{}:
		for _, item := range v {
			if itemStr, ok := item.(string); ok {
				if err := json.Unmarshal([]byte(itemStr), &innerPayloadData); err != nil {
					continue
				} else {
					break
				}
			}
		}
		if innerPayloadData == nil {
			return nil, fmt.Errorf("payload[1] is an array but no string element containing JSON object was found")
		}
	default:
		return nil, fmt.Errorf("payload[1] is not a string, object, or array, got %T", v)
	}

	return innerPayloadData, nil
}

func (c *VKVideoClient) parseNextFromAndSessionID(innerPayloadData map[string]interface{}) (string, int, error) {
	nextFrom, ok := innerPayloadData["nextFrom"].(string)
	if !ok {
		nextFrom = ""
	}

	var showcaseSessionID int
	if showcaseSessionIDVal, ok := innerPayloadData["showcaseSessionId"].(float64); ok {
		showcaseSessionID = int(showcaseSessionIDVal)
	} else {
		showcaseSessionID = 0
	}

	return nextFrom, showcaseSessionID, nil
}

func (c *VKVideoClient) extractPlaylists(innerPayloadData map[string]interface{}) ([][]interface{}, error) {
	playlistsRaw, ok := innerPayloadData["playlists"]
	if !ok {
		return nil, fmt.Errorf("playlists not found in inner payload data")
	}

	playlists, ok := playlistsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("playlists in inner payload data is not an array, got %T", playlistsRaw)
	}

	var result [][]interface{}
	for _, item := range playlists {
		playlist, ok := item.([]interface{})
		if !ok {
			continue
		}
		result = append(result, playlist)
	}

	return result, nil
}

func (c *VKVideoClient) FetchAllMoviesWithPagination(ctx context.Context, maxPages int) ([][]interface{}, error) {
	var allPlaylists [][]interface{}

	nextFrom := ""
	showcaseSessionID := 0

	pageCount := 0

	for {
		if maxPages > 0 && pageCount >= maxPages { // maxPages == 0 означает "без ограничений"
			fmt.Println("Reached maximum number of pages to fetch, stopping.")
			break
		}

		fmt.Printf("Fetching page %d, nextFrom: '%s', showcaseSessionId: %d\n", pageCount+1, nextFrom, showcaseSessionID)

		data, err := c.fetchNextPage(ctx, nextFrom, showcaseSessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", pageCount+1, err)
		}

		innerPayloadData, err := c.extractInnerPayloadData(data)
		if err != nil {
			return nil, fmt.Errorf("failed to extract inner payload data on page %d: %w", pageCount+1, err)
		}

		playlists, err := c.extractPlaylists(innerPayloadData)
		if err != nil {
			return nil, fmt.Errorf("failed to extract playlists on page %d: %w", pageCount+1, err)
		}

		allPlaylists = append(allPlaylists, playlists...)

		newNextFrom, newShowcaseSessionID, err := c.parseNextFromAndSessionID(innerPayloadData)
		if err != nil {
			fmt.Printf("Error parsing pagination info from page %d, stopping: %v\n", pageCount+1, err)
			break
		}

		if newNextFrom == "" {
			fmt.Println("No more pages (nextFrom is empty), stopping.")
			break
		}

		nextFrom = newNextFrom
		if newShowcaseSessionID != 0 {
			showcaseSessionID = newShowcaseSessionID
		}

		pageCount++
	}

	if len(allPlaylists) == 0 {
		return nil, fmt.Errorf("no playlists were fetched during pagination")
	}

	return allPlaylists, nil
}
