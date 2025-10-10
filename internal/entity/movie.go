package entity

import "time"

type Movie struct {
	ID           int
	Title        string
	Description  string
	AverageScore float64
	ReleaseDate  time.Time
}
