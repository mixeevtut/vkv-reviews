package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vkvideo-reviews/internal/api/http/v2"
	handlersV2 "vkvideo-reviews/internal/api/http/v2/handlers"
	"vkvideo-reviews/internal/cron"
	"vkvideo-reviews/internal/repository"
	"vkvideo-reviews/internal/usecase"
	"vkvideo-reviews/pkg/postgres"

	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Применить миграции
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Transactor
	transactor := postgres.NewTransactor(db)

	// Repository
	reviewRepo := repository.NewReviewsRepository(db)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMoviesRepository(db)

	// Service
	reviews := usecase.NewReviewUseCase(reviewRepo, transactor)
	user := usecase.NewUserUseCase(userRepo)
	movies := usecase.NewMoviesUseCase(movieRepo)

	// Импорт фильмов
	vkVideoParser := &repository.VKVideoParser{}
	movieImportUseCase := usecase.NewMovieImportUseCase(vkVideoParser, movieRepo)

	// Запускаем cron-задачу
	fetcher := handlersV2.NewVKVideoFetcher("https://vkvideo.ru/movies")
	cronJob := cron.NewMovieImportCron(movieImportUseCase, fetcher, 24*time.Hour)

	go cronJob.Start(ctx)
	// Handler
	// reviewHandler := handlersV1.NewReviewHandler(reviews)
	// userHandler := handlersV1.NewUserHandler(user)
	// moviesHandler := handlersV1.NewMoviesHandler(movies)

	// Handler v2
	reviewHandlerV2 := handlersV2.NewReviewHandler(reviews)
	moviesHandlerV2 := handlersV2.NewMoviesHandler(movies)
	userHandlerV2 := handlersV2.NewUserHandler(user)

	v2.NewRouter(
		reviewHandlerV2,
		moviesHandlerV2,
		userHandlerV2,
	).StartRouter()

	// Запускаем HTTP-сервер в горутине
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	// v1.StartRouter(
	//	reviewHandler,
	//	moviesHandler,
	//	userHandler,
	// )

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ждём сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel() // останавливаем cron

	// Останавливаем HTTP-сервер
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server exited")
}
