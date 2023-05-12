package data

import "greenlight.wook.net/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilter(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "0 보다 커야만 합니다")
	v.Check(f.Page <= 10_000_000, "page", "최대 1000만 까지 요청할 수 있습니다")
	v.Check(f.PageSize > 0, "page_size", "0 보다 커야만 합니다")
	v.Check(f.PageSize <= 100, "page_size", "최대 100 까지 요청할 수 있습니다")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "잘못된 sort 값 입니다")
}
