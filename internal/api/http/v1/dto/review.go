package dto

import (
	"time"
	"vkvideo-reviews/internal/entity"
)

type CreateReviewRequest struct {
	UserID  int
	MovieID int
	Score   int
	Content string
}

type ReviewDTO struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MovieID   int       `json:"movie_id"`
	Score     int       `json:"score"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToReviewResponse(r *entity.Review) ReviewDTO {
	return ReviewDTO{
		ID:        r.ID,
		UserID:    r.UserID,
		MovieID:   r.MovieID,
		Score:     r.Score,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type ErrorReview struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
