package cron

import (
	"context"
	"log"
	"time"

	"vkvideo-reviews/internal/api/http/v2/handlers"
	"vkvideo-reviews/internal/usecase"
)

type MovieImportCron struct {
	importer usecase.MovieImportUseCase
	fetcher  handlers.MovieDataFetcher
	interval time.Duration
}

func NewMovieImportCron(
	importer usecase.MovieImportUseCase,
	fetcher handlers.MovieDataFetcher,
	interval time.Duration,
) *MovieImportCron {
	return &MovieImportCron{
		importer: importer,
		fetcher:  fetcher,
		interval: interval,
	}
}

func (c *MovieImportCron) Start(ctx context.Context) {
	log.Println("Starting movie import cron...")

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// Запускаем сразу при старте (опционально)
	c.runOnce(ctx)

	for {
		select {
		case <-ticker.C:
			c.runOnce(ctx)
		case <-ctx.Done():
			log.Println("Stopping movie import cron")
			return
		}
	}
}

func (c *MovieImportCron) runOnce(ctx context.Context) {
	log.Println("Fetching and importing movies...")

	data, err := c.fetcher.FetchMoviesData(ctx)
	if err != nil {
		log.Printf("Failed to fetch movie data: %v", err)
		return
	}

	result, err := c.importer.ImportMoviesFromJSON(ctx, data)
	if err != nil {
		log.Printf("Failed to import movies: %v", err)
		return
	}

	log.Printf("Imported %d movies (total: %d)", result.Imported, result.Total)
}
