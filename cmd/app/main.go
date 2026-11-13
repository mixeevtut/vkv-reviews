package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"vkvideo-reviews/internal/api/http/v2"
	handlersV2 "vkvideo-reviews/internal/api/http/v2/handlers"
	"vkvideo-reviews/internal/config"
	"vkvideo-reviews/internal/cron"
	"vkvideo-reviews/internal/fetcher"
	"vkvideo-reviews/internal/repository"
	"vkvideo-reviews/internal/usecase"
	"vkvideo-reviews/pkg/postgres"
	"vkvideo-reviews/pkg/redis"

	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	//Redis
	rds, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatalf("redis init: %v", err)
	}
	defer func() { _ = rds.Close() }()
	if err := rds.EnsureGroup(ctx); err != nil {
		log.Fatalf("initialize Redis stream group: %v", err)
	}

	// Repository
	reviewRepo := repository.NewReviewsRepository(db)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMoviesRepository(db)
	eventRepo := repository.NewRecalcEventRepo(db)

	// Service
	reviews := usecase.NewReviewUseCase(reviewRepo, eventRepo, transactor)
	user := usecase.NewUserUseCase(userRepo)
	movies := usecase.NewMoviesUseCase(movieRepo, reviewRepo)
	ratingRecalc := usecase.NewRatingRecalcUseCase(reviewRepo, movieRepo, transactor)

	producer := usecase.NewProducer(eventRepo, rds)
	consumer := usecase.NewConsumer(rds, *ratingRecalc)

	vkVideoClient := fetcher.NewVKVideoClient("https://vkvideo.ru/movies_serials/films")

	// Импорт фильмов
	vkVideoParser := &repository.VKVideoParser{}
	movieImportUseCase := usecase.NewMovieImportUseCase(vkVideoClient, vkVideoParser, movieRepo)

	// Запускаем cron-задачу (передаем maxPages)
	maxPagesPerRun := 10
	cronJob := cron.NewMovieImportCron(movieImportUseCase, 24*time.Hour, maxPagesPerRun)

	// Handler v2
	reviewHandlerV2 := handlersV2.NewReviewHandler(reviews)
	moviesHandlerV2 := handlersV2.NewMoviesHandler(movies)
	userHandlerV2 := handlersV2.NewUserHandler(user)

	handler, err := v2.NewRouter(reviewHandlerV2, moviesHandlerV2, userHandlerV2).Handler()
	if err != nil {
		log.Fatalf("Failed to build v2 router: %v", err)
	}

	var wg sync.WaitGroup
	wg.Go(func() { consumer.Run(ctx) })
	wg.Go(func() { producer.Run(ctx, 2*time.Second, 100) })
	wg.Go(func() { cronJob.Start(ctx) })

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	wg.Wait()

	log.Println("Server exited")
}
