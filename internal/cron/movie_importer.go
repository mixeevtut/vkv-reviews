package cron

import (
	"context"
	"log"
	"time"

	"vkvideo-reviews/internal/usecase"
)

type MovieImportCron struct {
	importer usecase.MovieImportUseCase
	interval time.Duration
	maxPages int
}

func NewMovieImportCron(
	importer usecase.MovieImportUseCase,
	interval time.Duration,
	maxPages int,
) *MovieImportCron {
	return &MovieImportCron{
		importer: importer,
		interval: interval,
		maxPages: maxPages,
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

	result, err := c.importer.ImportMovies(ctx, c.maxPages)
	if err != nil {
		log.Printf("Failed to import movies: %v", err)
		return
	}

	log.Printf("Imported %d movies (total processed: %d)", result.Imported, result.Total)
}
