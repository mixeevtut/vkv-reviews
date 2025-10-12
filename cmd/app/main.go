package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	_ "vkvideo-reviews/docs"
	v1 "vkvideo-reviews/internal/api/http/v1"
	"vkvideo-reviews/internal/api/http/v1/handlers"
	"vkvideo-reviews/internal/repository"
	"vkvideo-reviews/internal/usecase"
)

// @title           VKVideo-Reviews API
// @description     API для  рецензий
// @BasePath        /api/v1
// @schemes         http

func main() {

	ctx := context.Background()

	// Получаем URL базы данных из переменной окружения
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/tasks_db?sslmode=disable"
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// repository
	reviewRepo := repository.NewReviewsRepository(db)
	movieRepo := repository.NewMoviesRepository(db)

	// Service
	reviews := usecase.NewReviewUseCase(reviewRepo)
	movies := usecase.NewGetMovieUseCase(movieRepo)

	// Handler
	reviewHandler := handlers.NewReviewHandler(reviews)
	moviesHandler := handlers.NewMoviesHandler(movies)

	v1.StartRouter(
		reviewHandler,
		moviesHandler,
	)

}
