package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"owenHochwald.greenlight/internal/validator"
)

var ErrMovieNotFound = errors.New("movie not found")
var ErrMovieEditConflict = errors.New("edit conflict")

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

// mock movie model methods
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

func (m MovieMockModel) GetAll() (*[]Movie, error) {
	return nil, nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, nil
	}
	var movie Movie
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1`

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreateAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		return nil, ErrMovieNotFound
	}

	return &movie, nil
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	err := m.DB.QueryRow(query, args...).Scan(
		&movie.ID,
		&movie.CreateAt,
		&movie.Version,
	)
	return err
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
            UPDATE movies
            SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
            WHERE id = $5 AND version = $6
            RETURNING version`

	var args = []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}
	_, err := m.DB.Exec(query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrMovieEditConflict

		default:
			return ErrMovieNotFound
		}
	}

	return nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return nil
	}
	query := `
         DELETE FROM movies
         WHERE id = $1`

	res, err := m.DB.Exec(query, id)

	if err != nil {
		return err
	}

	numberRows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	fmt.Println(numberRows)
	if numberRows == 0 {
		return ErrMovieNotFound
	}

	return nil
}

func (m MovieModel) GetAll() (*[]Movie, error) {
	var movies []Movie
	query := `
			SELECT * FROM movies`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.CreateAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return &movies, nil
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
