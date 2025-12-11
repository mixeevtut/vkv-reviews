package dto

import (
	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/entity"
)

func ToMovieResponse(movie entity.Movie) apigen.Movie {
	return apigen.Movie{
		Id:           movie.ID,
		Title:        movie.Title,
		Description:  movie.Description,
		AverageScore: movie.AverageScore,
		ReleaseDate:  movie.ReleaseDate,
		Reviews:      movie.Reviews,
	}
}
