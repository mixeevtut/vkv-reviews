package repository

import (
	"context"
	"database/sql"
	"errors"

	"vkvideo-reviews/internal/entity"
)

type MoviesRepository interface {
	Create(ctx context.Context, movie *entity.Movie) error
	CreateBatch(ctx context.Context, movies []*entity.Movie) error
	GetByID(ctx context.Context, ID int) (*entity.Movie, error)
	GetAll(ctx context.Context) ([]*entity.Movie, error)
	Exists(ctx context.Context, ID int) (bool, error)
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

func (r *Movies) Create(ctx context.Context, movie *entity.Movie) error {
	if movie == nil {
		return errors.New("movie is nil")
	}

	query := `
		INSERT INTO movies (id, title, description, average_score, release_date)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			average_score = EXCLUDED.average_score,
			release_date = EXCLUDED.release_date,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, query,
		movie.ID,
		movie.Title,
		movie.Description,
		movie.AverageScore,
		movie.ReleaseDate,
	)

	return err
}

func (r *Movies) CreateBatch(ctx context.Context, movies []*entity.Movie) error {
	if len(movies) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO movies (id, title, description, average_score, release_date)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			average_score = EXCLUDED.average_score,
			release_date = EXCLUDED.release_date,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, movie := range movies {
		if movie == nil {
			continue
		}

		_, err := stmt.ExecContext(ctx,
			movie.ID,
			movie.Title,
			movie.Description,
			movie.AverageScore,
			movie.ReleaseDate,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Movies) Exists(ctx context.Context, ID int) (bool, error) {
	if ID <= 0 {
		return false, errors.New("invalid movie ID")
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)`

	err := r.db.QueryRowContext(ctx, query, ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
