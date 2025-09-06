package data

import "database/sql"

type Models struct {
	Movies interface {
		Get(id int64) (*Movie, error)
		Insert(movie *Movie) error
		Update(movie *Movie) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Movie, error)
	}
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Movies: &MovieModel{DB: db},
	}
}

func newMockModels() *Models {
	return &Models{
		Movies: &MovieMockModel{},
	}
}
