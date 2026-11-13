package usecase

import (
	"context"
	"fmt"
	"log"

	"vkvideo-reviews/internal/fetcher"
	"vkvideo-reviews/internal/repository"
)

type MovieImportUseCase interface {
	ImportMovies(ctx context.Context, maxPages int) (*ImportResult, error)
}

type MovieImport struct {
	client    *fetcher.VKVideoClient
	parser    *repository.VKVideoParser
	movieRepo repository.MoviesRepository
}

func NewMovieImportUseCase(client *fetcher.VKVideoClient, parser *repository.VKVideoParser, movieRepo repository.MoviesRepository) *MovieImport {
	return &MovieImport{
		client:    client,
		parser:    parser,
		movieRepo: movieRepo,
	}
}

func (uc *MovieImport) ImportMovies(ctx context.Context, maxPages int) (*ImportResult, error) {
	items, err := uc.client.FetchAllMoviesWithPagination(ctx, maxPages)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies from VK: %w", err)
	}
	movies, err := uc.parser.ParseMovies(items)
	if err != nil {
		return nil, fmt.Errorf("failed to parse movies from VK: %w", err)
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
