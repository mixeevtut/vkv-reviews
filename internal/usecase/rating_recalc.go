package usecase

import (
	"context"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/internal/repository"
	"vkvideo-reviews/pkg/transaction"
)

type RatingRecalculator interface {
	RecalcReview(ctx context.Context, movieID int) error
}

type RatingRecalcUseCase struct {
	movieRepo  repository.MoviesRepository
	reviewRepo repository.ReviewRepository
	transactor transaction.Transactor
}

func NewRatingRecalcUseCase(reviewRepo repository.ReviewRepository, movies repository.MoviesRepository, transactor transaction.Transactor) *RatingRecalcUseCase {
	return &RatingRecalcUseCase{
		movieRepo:  movies,
		reviewRepo: reviewRepo,
		transactor: transactor,
	}
}

func (r *RatingRecalcUseCase) RecalcReview(ctx context.Context, movieID int) error {
	err := r.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		reviews, err := r.reviewRepo.GetReviewsByMovie(ctx, movieID)
		if err != nil {
			return err
		}
		avg := averageRating(reviews)
		err = r.movieRepo.UpdateAvgRating(ctx, movieID, avg)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func averageRating(reviews []entity.Review) float64 {
	if len(reviews) == 0 {
		return 0
	}

	var sum float64
	for _, review := range reviews {
		sum += float64(review.Score)
	}
	return sum / float64(len(reviews))
}
