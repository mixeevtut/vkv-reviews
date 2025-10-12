package usecase

import (
	"context"
	"vkvideo-reviews/internal/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u entity.User) error
}

type UserUseCase struct {
	userRepo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

func (s UserUseCase) CreateUser(ctx context.Context, u entity.User) error {
	return s.userRepo.CreateUser(ctx, u)
}
