package data

import (
	"database/sql"
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

type MovieModel struct {
	DB *sql.DB
}

type MovieMockModel struct{}

func (m MovieMockModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieMockModel) Insert(movie *Movie) error {
	return nil
}

func (m MovieMockModel) Update(movie *Movie) error {
	return nil
}

func (m MovieMockModel) Delete(id int64) error {
	return nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieModel) Insert(movie *Movie) error {
	return nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
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
