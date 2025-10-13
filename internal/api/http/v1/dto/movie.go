package dto

import (
	"time"
	"vkvideo-reviews/internal/entity"
)

type ModelMovie struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	AverageScore float64   `json:"average_score"`
	ReleaseDate  time.Time `json:"release_date"`
}

type MoviesResponse struct {
	Movies []*ModelMovie `json:"movies"`
}

func (m *ModelMovie) ToMovie() entity.Movie {
	return entity.Movie{
		ID:           m.ID,
		Title:        m.Title,
		Description:  m.Description,
		AverageScore: m.AverageScore,
		ReleaseDate:  m.ReleaseDate,
	}
}

func FromMovie(movie entity.Movie) *ModelMovie {
	return &ModelMovie{
		ID:           movie.ID,
		Title:        movie.Title,
		Description:  movie.Description,
		AverageScore: movie.AverageScore,
		ReleaseDate:  movie.ReleaseDate,
	}
}
