package repository

import (
	"context"
	"database/sql"
	"errors"
	"vkvideo-reviews/internal/entity"
)

type MoviesRepository interface {
	GetByID(ctx context.Context, ID int) (*entity.Movie, error)
}

type Movies struct {
	db *sql.DB
}

func NewMoviesRepository(db *sql.DB) *Movies {
	return &Movies{db: db}
}

func (r *Movies) GetByID(ctx context.Context, ID int) (*entity.Movie, error) {
	if ID <= 0 {
		return nil, errors.New("invalid movie ID")
	}

	var movie entity.Movie
	query := `SELECT id, title, description, avgscore, releasedate FROM movies WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, ID).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.AverageScore,
		&movie.ReleaseDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}

	return &movie, nil
}
