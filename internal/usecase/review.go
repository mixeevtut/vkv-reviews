package usecase

import (
	"context"
	"github.com/pkg/errors"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/internal/repository"
	"vkvideo-reviews/pkg/transaction"
)

type ReviewUseCase struct {
	reviewRepo repository.ReviewRepository
	event      repository.RecalcEventRepository
	transactor transaction.Transactor
}

func NewReviewUseCase(repo repository.ReviewRepository, eventRepo repository.RecalcEventRepository, transactor transaction.Transactor) *ReviewUseCase {
	return &ReviewUseCase{
		reviewRepo: repo,
		event:      eventRepo,
		transactor: transactor,
	}
}

func (r *ReviewUseCase) GetReview(ctx context.Context, ID int) (*entity.Review, error) {
	review, err := r.reviewRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewUseCase) CreateReview(ctx context.Context, rw entity.Review) (int, error) {
	var result int
	err := r.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		reviewID, err := r.reviewRepo.Create(ctx, rw)
		if err != nil {
			return err
		}
		result = reviewID
		err = r.event.Add(ctx, rw.MovieID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return -1, errors.Wrap(err, "review creation failed")
	}

	return result, nil
}

func (r *ReviewUseCase) UpdateReview(ctx context.Context, rw entity.ReviewUpdate) (*entity.Review, error) {
	var result *entity.Review
	err := r.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		review, err := r.reviewRepo.Update(ctx, rw)
		if err != nil {
			return err
		}
		result = review
		err = r.event.Add(ctx, review.MovieID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "error executing update review")
	}

	return result, nil
}

func (r *ReviewUseCase) DeleteReview(ctx context.Context, ID int) error {
	err := r.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		review, err := r.reviewRepo.GetByID(ctx, ID)
		if err != nil {
			return err
		}
		movieId := review.MovieID
		err = r.reviewRepo.DeleteByID(ctx, ID)
		if err != nil {
			return err
		}
		err = r.event.Add(ctx, movieId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error executing delete review")
	}
	return nil
}
