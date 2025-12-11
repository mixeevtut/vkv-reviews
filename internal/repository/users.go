package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"vkvideo-reviews/internal/entity"
)

type UserRepo struct {
	DB *sql.DB
}

type user struct {
	Id    int
	Name  string
	Email string
}

func (u *user) toEntity() *entity.User {
	return &entity.User{
		ID:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, u entity.User) (*entity.User, error) {
	query := `INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id, name, email`
	var usr user

	scan := r.DB.QueryRowContext(ctx, query, u.Name, u.Email)
	err := scan.Scan(&usr.Id, &usr.Name, &usr.Email)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("email already exists: %w", err)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return usr.toEntity(), nil
}
