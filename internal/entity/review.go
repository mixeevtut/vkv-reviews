package entity

import "time"

type Review struct {
	ID        int
	UserID    int
	MovieID   int
	Score     int
	Content   string
	CreatedAt time.Time
}
