package entity

import "errors"

var (
	ErrMovieNotFound  = errors.New("movie not found")
	ErrReviewNotFound = errors.New("review not found")
)
