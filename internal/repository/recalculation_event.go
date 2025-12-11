package repository

import (
	"context"
	"database/sql"
	"time"
	"vkvideo-reviews/internal/entity"
)

type RecalcEventRepository interface {
	FetchUnsent(ctx context.Context, limit int) ([]entity.RecalcEvent, error)
	MarkSent(ctx context.Context, id int64) error
	Add(ctx context.Context, movieID int) error
}

type RecalcEventRepo struct{ db *sql.DB }

func NewRecalcEventRepo(db *sql.DB) *RecalcEventRepo { return &RecalcEventRepo{db: db} }

func (r *RecalcEventRepo) FetchUnsent(ctx context.Context, limit int) ([]entity.RecalcEvent, error) {
	const q = `SELECT id, movie_id, created_at, sent_at FROM recalculation_event WHERE sent_at IS NULL ORDER BY id LIMIT $1`

	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, entity.ErrEventFetch
	}
	defer rows.Close()

	var res []entity.RecalcEvent
	for rows.Next() {
		var e entity.RecalcEvent
		if err := rows.Scan(&e.ID, &e.MovieID, &e.CreatedAt, &e.SentAt); err != nil {
			return nil, entity.ErrEventFetch
		}
		res = append(res, e)
	}
	return res, nil
}

func (r *RecalcEventRepo) MarkSent(ctx context.Context, id int64) error {
	const q = `UPDATE recalculation_event SET sent_at = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, q, id, time.Now().UTC())
	if err != nil {
		return entity.ErrEventMarkSent
	}
	return nil
}

func (r *RecalcEventRepo) Add(ctx context.Context, movieID int) error {
	const q = `INSERT INTO recalculation_event(movie_id) VALUES ($1)`

	_, err := r.db.ExecContext(ctx, q, movieID)
	if err != nil {
		return entity.ErrEventAdd
	}

	return nil
}
