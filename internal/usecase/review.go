package usecase

import (
	"context"
	"vkvideo-reviews/internal/entity"
)

type ReviewRepository interface {
	GetByID(ctx context.Context, ID int) (*entity.Review, error)
	Create(ctx context.Context, r entity.Review) (int, error)
	Update(ctx context.Context, r entity.Review) (*entity.Review, error)
	DeleteByID(ctx context.Context, ID int) error
}

type ReviewUseCase struct {
	reviewRepo ReviewRepository
}

func NewReviewUseCase(repo ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{
		reviewRepo: repo,
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
	reviewID, err := r.reviewRepo.Create(ctx, rw)
	if err != nil {
		return -1, err
	}

	return reviewID, nil
}

func (r *ReviewUseCase) UpdateReview(ctx context.Context, rw entity.Review) (*entity.Review, error) {

	review, err := r.reviewRepo.Update(ctx, rw)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewUseCase) DeleteReview(ctx context.Context, ID int) error {
	err := r.reviewRepo.DeleteByID(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
