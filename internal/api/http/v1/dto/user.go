package dto

import (
	"errors"
	"strings"
	"vkvideo-reviews/internal/entity"
)

type ModelUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (m *ModelUser) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(m.Email) == "" {
		return errors.New("email is required")
	}
	return nil
}

func (m *ModelUser) ToUser() entity.User {
	return entity.User{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}
}
