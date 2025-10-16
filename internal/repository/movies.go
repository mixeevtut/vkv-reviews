package repository

import (
	"context"
	"database/sql"
	"errors"
	"vkvideo-reviews/internal/entity"
)

type MoviesRepository interface {
	GetByID(ctx context.Context, ID int) (*entity.Movie, error)
	GetAll(ctx context.Context) ([]*entity.Movie, error)
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
	query := `SELECT id, title, description, average_score, release_date FROM movies WHERE id = $1`

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

func (r *Movies) GetAll(ctx context.Context) ([]*entity.Movie, error) {
	query := `SELECT id, title, description, average_score, release_date FROM movies`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*entity.Movie
	for rows.Next() {
		var movie entity.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.AverageScore,
			&movie.ReleaseDate,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}
