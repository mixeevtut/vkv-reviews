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

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, u entity.User) error {
	query := `INSERT INTO users(name, email) VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, u.Name, u.Email)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("email already exists: %w", err)
		}
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
