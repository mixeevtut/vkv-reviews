package entity

import "time"

type Review struct {
	ID        int
	UserID    int
	MovieID   int
	Score     int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewUpdate struct {
	ID      int
	Score   *int
	Content *string
}
