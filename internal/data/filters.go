package data

import (
	"slices"

	"owenHochwald.greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "Page is greater than 0")
	v.Check(f.Page <= 10_000_000, "page", "Page is less than 10 million")

	v.Check(f.PageSize > 0, "page_size", "Page size is greater than 0")
	v.Check(f.PageSize <= 20, "page_size", "Page size is less than 10")

	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "Sort option must be valid")
}

func (f Filters) sortDirection() string {
	if f.Sort != "" && f.Sort[0] == '-' {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) sortColumn() string {
	if slices.Contains(f.SortSafeList, f.Sort) {
		if f.Sort[0] == '-' {
			return f.Sort[1:]
		}
	}
	//panic("Unsafe sort value - Triggering the ultra protection fail safes")
	return "id"
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) limit() int {
	return f.PageSize
}
