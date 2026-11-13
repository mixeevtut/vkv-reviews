package repository

import (
	"fmt"
	"strconv"
	"strings"

	"vkvideo-reviews/internal/entity"
)

type VKVideoParser struct{}

func (p *VKVideoParser) ParseMovies(items [][]interface{}) ([]*entity.Movie, error) {
	var movies []*entity.Movie

	for _, item := range items {
		movie, err := p.parseItem(item)
		if err != nil {
			continue
		}
		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return nil, fmt.Errorf("no valid movies found in the catalog")
	}

	return movies, nil
}

func (p *VKVideoParser) parseItem(item []interface{}) (*entity.Movie, error) {
	if len(item) > 20 {
		if path, ok := item[20].(string); ok && strings.HasPrefix(path, "/video") {
			id, err := p.extractVideoID(strings.TrimPrefix(path, "/video"))
			title, hasTitle := item[3].(string)
			if err == nil && id > 0 && hasTitle && title != "" {
				return &entity.Movie{ID: id, Title: title}, nil
			}
		}
	}
	return p.parsePlaylistItem(item)
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
