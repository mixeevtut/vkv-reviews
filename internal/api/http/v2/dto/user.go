package dto

import (
	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/entity"
)

func ToUserEntity(u apigen.UserCreate) entity.User {
	return entity.User{
		Name:  u.Name,
		Email: string(u.Email),
	}
}

func ToUserCreateResponse(id int) apigen.UserCreateResponse {
	return apigen.UserCreateResponse{
		Id: id,
	}
}
