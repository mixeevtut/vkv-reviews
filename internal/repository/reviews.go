package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"time"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/pkg/postgres"
)

type ReviewRepoPG struct {
	db *sql.DB
}

func NewReviewsRepository(db *sql.DB) *ReviewRepoPG {
	return &ReviewRepoPG{
		db: db,
	}
}

type row struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	MovieID   int
	Score     int
	Content   string
}

func (r *row) toEntity() *entity.Review {
	return &entity.Review{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		UserID:    r.UserID,
		MovieID:   r.MovieID,
		Score:     r.Score,
		Content:   r.Content,
	}
}

type Executor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (p *ReviewRepoPG) getExecutor(ctx context.Context) Executor {
	if tx, ok := postgres.GetTx(ctx); ok {
		return tx
	}
	logger.Warnf("transaction disabled")
	return p.db
}

func (p *ReviewRepoPG) GetByID(ctx context.Context, ID int) (*entity.Review, error) {
	const q = `SELECT id, created_at, updated_at, user_id, movie_id, score, content FROM reviews WHERE id=$1 AND is_deleted = FALSE`
	var rw row

	exec := p.getExecutor(ctx)

	scan := exec.QueryRowContext(ctx, q, ID)
	err := scan.Scan(
		&rw.ID,
		&rw.CreatedAt,
		&rw.UpdatedAt,
		&rw.UserID,
		&rw.MovieID,
		&rw.Score,
		&rw.Content,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrReviewNotFound
		}
		return nil, err
	}

	return rw.toEntity(), nil
}

func (p *ReviewRepoPG) Create(ctx context.Context, r entity.Review) (int, error) {
	const q = `INSERT INTO reviews(created_at, updated_at, user_id, movie_id, score, content) VALUES (now(), now(), $1, $2, $3, $4) RETURNING id`
	var rw row

	exec := p.getExecutor(ctx)
	scan := exec.QueryRowContext(ctx, q, r.UserID, r.MovieID, r.Score, r.Content)
	err := scan.Scan(&rw.ID)
	if err != nil {
		return -1, err
	}

	return rw.ID, nil
}

func (p *ReviewRepoPG) Update(ctx context.Context, r entity.Review) (*entity.Review, error) {
	const q = `UPDATE reviews SET score = $1, content = $2, updated_at = now() WHERE id = $3 AND is_deleted = FALSE RETURNING id, created_at, updated_at, user_id, movie_id, score, content;`
	var rw row

	exec := p.getExecutor(ctx)
	scan := exec.QueryRowContext(ctx, q, r.Score, r.Content, r.ID)
	err := scan.Scan(&rw.ID, &rw.CreatedAt, &rw.UpdatedAt, &rw.UserID, &rw.MovieID, &rw.Score, &rw.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrReviewNotFound
		}
		return nil, err
	}

	return rw.toEntity(), nil
}

func (p *ReviewRepoPG) DeleteByID(ctx context.Context, ID int) error {
	const q = `UPDATE reviews SET is_deleted = TRUE WHERE id = $1 AND is_deleted = FALSE RETURNING id`
	var rw row

	exec := p.getExecutor(ctx)
	scan := exec.QueryRowContext(ctx, q, ID)
	err := scan.Scan(&rw.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrReviewNotFound
		}
		return err
	}

	return nil
}
