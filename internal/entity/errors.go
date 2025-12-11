package entity

import "errors"

var (
	ErrMovieNotFound  = errors.New("movie not found")
	ErrReviewNotFound = errors.New("review not found")
)

var (
	ErrEventFetch    = errors.New("event fetch unsent")
	ErrEventMarkSent = errors.New("event mark sent")
	ErrEventAdd      = errors.New("event add")
)
