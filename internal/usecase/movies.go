package usecase

import (
	"context"
	"fmt"

	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/internal/repository"
)

type MoviesUseCase interface {
	GetMovie(ctx context.Context, ID int) (*entity.Movie, error)
	GetMovies(ctx context.Context) ([]*entity.Movie, error)
}

type Movies struct {
	movieRepo repository.MoviesRepository
}

func NewMoviesUseCase(movieRepo repository.MoviesRepository) *Movies {
	return &Movies{
		movieRepo: movieRepo,
	}
}

func (m *Movies) GetMovie(ctx context.Context, ID int) (*entity.Movie, error) {
	movie, err := m.movieRepo.GetByID(ctx, ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return movie, nil
}

func (m *Movies) GetMovies(ctx context.Context) ([]*entity.Movie, error) {
	movies, err := m.movieRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	return movies, nil
}
