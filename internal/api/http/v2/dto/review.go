package dto

import (
	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/entity"
)

func CreateReviewToEntity(a apigen.ReviewCreate) entity.Review {
	return entity.Review{
		UserID:  a.UserId,
		MovieID: a.MovieId,
		Score:   a.Score,
		Content: a.Content,
	}
}

func UpdateReviewToEntity(id int, a apigen.ReviewUpdate) entity.Review {
	return entity.Review{
		ID:      id,
		Score:   *a.Score,
		Content: *a.Content,
	}
}

func ToGetReviewResponse(r *entity.Review) apigen.Review {
	return apigen.Review{
		Id:        r.ID,
		UserId:    r.UserID,
		MovieId:   r.MovieID,
		Content:   r.Content,
		Score:     r.Score,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func ToCreateReviewResponse(id int) apigen.ReviewCreateResponse {
	return apigen.ReviewCreateResponse{
		Id: id,
	}
}

func ToUpdateReviewResponse(r *entity.Review) apigen.Review {
	return apigen.Review{
		Id:        r.ID,
		UserId:    r.UserID,
		MovieId:   r.MovieID,
		Content:   r.Content,
		Score:     r.Score,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
