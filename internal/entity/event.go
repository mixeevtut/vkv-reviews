package entity

import "time"

type RecalcEvent struct {
	ID        int64
	MovieID   int64
	CreatedAt time.Time
	SentAt    *time.Time
}
