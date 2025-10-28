package repository

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"vkvideo-reviews/internal/entity"
)

type VKVideoParser struct{}

func (p *VKVideoParser) ParseMoviesJSON(jsonData []byte) ([]*entity.Movie, error) {
	var response struct {
		Payload []interface{} `json:"payload"`
	}

	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal root JSON: %w", err)
	}

	if len(response.Payload) < 2 {
		return nil, fmt.Errorf("invalid payload structure: expected at least 2 elements")
	}

	innerPayload, ok := response.Payload[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("payload[1] is not an array")
	}

	if len(innerPayload) < 1 {
		return nil, fmt.Errorf("inner payload is empty")
	}

	jsonStr, ok := innerPayload[0].(string)
	if !ok {
		return nil, fmt.Errorf("inner payload[0] is not a string")
	}

	var videoData struct {
		Playlists [][]interface{} `json:"playlists"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &videoData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner JSON: %w", err)
	}

	return p.parsePlaylists(videoData.Playlists)
}

func (p *VKVideoParser) parsePlaylists(playlists [][]interface{}) ([]*entity.Movie, error) {
	var movies []*entity.Movie

	for _, playlist := range playlists {
		movie, err := p.parsePlaylistItem(playlist)
		if err != nil {
			continue
		}
		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return nil, fmt.Errorf("no valid movies found in playlists")
	}

	return movies, nil
}

func (p *VKVideoParser) parsePlaylistItem(playlist []interface{}) (*entity.Movie, error) {
	if len(playlist) < 8 {
		return nil, fmt.Errorf("invalid playlist item length: %d", len(playlist))
	}

	movie := &entity.Movie{
		AverageScore: 0.0,
	}

	if title, ok := playlist[0].(string); ok {
		movie.Title = title
	} else {
		return nil, fmt.Errorf("invalid title")
	}

	if videoID, ok := playlist[7].(string); ok {
		id, err := p.extractVideoID(videoID)
		if err != nil {
			return nil, fmt.Errorf("failed to extract video ID: %w", err)
		}
		movie.ID = id
	} else {
		return nil, fmt.Errorf("invalid video ID")
	}

	return movie, nil
}

func (p *VKVideoParser) extractVideoID(videoIDStr string) (int, error) {
	parts := strings.Split(videoIDStr, "_")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid video ID format: %s", videoIDStr)
	}
	return strconv.Atoi(parts[1])
}
