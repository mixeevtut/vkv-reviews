package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"vkvideo-reviews/internal/api/http/v2"
	handlersV2 "vkvideo-reviews/internal/api/http/v2/handlers"
	"vkvideo-reviews/internal/repository"
	"vkvideo-reviews/internal/usecase"
	"vkvideo-reviews/pkg/postgres"

	_ "github.com/lib/pq"
)

func main() {

	ctx := context.Background()

	// Получаем URL базы данных из переменной окружения
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/vkvideo_reviews_db?sslmode=disable"
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if errBDClose := db.Close(); errBDClose != nil {
			log.Fatalf("Failed to close database connection: %v", errBDClose)
		}
	}()

	// Проверяем соединение
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	//Transactor
	transactor := postgres.NewTransactor(db)

	// repository
	reviewRepo := repository.NewReviewsRepository(db)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMoviesRepository(db)

	// Service
	reviews := usecase.NewReviewUseCase(reviewRepo, transactor)
	user := usecase.NewUserUseCase(userRepo)
	movies := usecase.NewGetMovieUseCase(movieRepo)

	// Handler
	//reviewHandler := handlersV1.NewReviewHandler(reviews)
	//userHandler := handlersV1.NewUserHandler(user)
	//moviesHandler := handlersV1.NewMoviesHandler(movies)

	// Handler v2
	reviewHandlerV2 := handlersV2.NewReviewHandler(reviews)
	moviesHandlerV2 := handlersV2.NewMoviesHandler(movies)
	userHandlerV2 := handlersV2.NewUserHandler(user)

	v2.NewRouter(
		reviewHandlerV2,
		moviesHandlerV2,
		userHandlerV2,
	).StartRouter()

	//v1.StartRouter(
	//	reviewHandler,
	//	moviesHandler,
	//	userHandler,
	//)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
