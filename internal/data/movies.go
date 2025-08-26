package data

import (
	"time"

	"owenHochwald.greenlight/internal/validator"
)

type Movie struct {
	ID       int64     `json:"id"`
	CreateAt time.Time `json:"-"` // '-' hides the field form JSON
	Title    string    `json:"title"`
	Year     int32     `json:"year"`
	Runtime  Runtime   `json:"runtime"`
	Genres   []string  `json:"genres"`
	Version  int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "Title is required")
	v.Check(len(movie.Title) < 500, "title", "Must not be longer than 500 characters")

	v.Check(movie.Year > 1888, "year", "Must be a valid year")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "Must be a valid year")
	v.Check(movie.Year != 0, "year", "Must be a valid year")

	v.Check(movie.Runtime != 0, "runtime", "Runtime must be provided")
	v.Check(movie.Runtime > 0, "runtime", "Must be a valid runtime")

	v.Check(len(movie.Genres) > 0, "genres", "Must provide at least one genre")
	v.Check(len(movie.Genres) < 5, "genres", "Must not be longer than 5 genres")

	v.Check(validator.Unique(movie.Genres), "genres", "Must be unique")
}
