package usecase

import (
	"context"
	"fmt"
	"vkvideo-reviews/internal/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u entity.User) (*entity.User, error)
}

type UserUseCase struct {
	userRepo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

func (s UserUseCase) CreateUser(ctx context.Context, u entity.User) (int, error) {
	user, err := s.userRepo.CreateUser(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return user.ID, nil
}
