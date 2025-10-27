package usecase

import (
	"context"
	"github.com/pkg/errors"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/pkg/transaction"
)

type ReviewRepository interface {
	GetByID(ctx context.Context, ID int) (*entity.Review, error)
	Create(ctx context.Context, r entity.Review) (int, error)
	Update(ctx context.Context, r entity.Review) (*entity.Review, error)
	DeleteByID(ctx context.Context, ID int) error
}

type ReviewUseCase struct {
	reviewRepo ReviewRepository
	transactor transaction.Transactor
}

func NewReviewUseCase(repo ReviewRepository, transactor transaction.Transactor) *ReviewUseCase {
	return &ReviewUseCase{
		reviewRepo: repo,
		transactor: transactor,
	}
}

func (r *ReviewUseCase) GetReview(ctx context.Context, ID int) (*entity.Review, error) {
	review, err := r.reviewRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, errGetReview
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

		return nil
	})

	if err != nil {
		return -1, errors.Wrap(err, "review creation failed")
	}

	return result, nil
}

func (r *ReviewUseCase) UpdateReview(ctx context.Context, rw entity.Review) (*entity.Review, error) {
	var result *entity.Review
	err := r.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		review, err := r.reviewRepo.Update(ctx, rw)
		if err != nil {
			return err
		}
		result = review

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "error executing update review")
	}

	return result, nil
}

func (r *ReviewUseCase) DeleteReview(ctx context.Context, ID int) error {
	err := r.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err := r.reviewRepo.DeleteByID(ctx, ID)
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
