package usecase

import (
	"context"
	"fmt"
	"log"

	"vkvideo-reviews/internal/repository"
)

type MovieImportUseCase interface {
	ImportMoviesFromJSON(ctx context.Context, jsonData []byte) (*ImportResult, error)
}

type MovieImport struct {
	parser    *repository.VKVideoParser
	movieRepo repository.MoviesRepository
}

func NewMovieImportUseCase(parser *repository.VKVideoParser, movieRepo repository.MoviesRepository) *MovieImport {
	return &MovieImport{
		parser:    parser,
		movieRepo: movieRepo,
	}
}

func (uc *MovieImport) ImportMoviesFromJSON(ctx context.Context, jsonData []byte) (*ImportResult, error) {
	movies, err := uc.parser.ParseMoviesJSON(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse movies from JSON: %w", err)
	}

	if len(movies) == 0 {
		return &ImportResult{Total: 0, Imported: 0}, nil
	}

	if err := uc.movieRepo.CreateBatch(ctx, movies); err != nil {
		return nil, fmt.Errorf("failed to save movies to database: %w", err)
	}

	log.Printf("Successfully imported %d movies", len(movies))

	return &ImportResult{
		Total:    len(movies),
		Imported: len(movies),
	}, nil
}

type ImportResult struct {
	Total    int
	Imported int
	Errors   []string
}
