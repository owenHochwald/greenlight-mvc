package data

import (
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
