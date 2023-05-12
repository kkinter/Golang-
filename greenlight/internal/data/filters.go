package data

import (
	"strings"

	"greenlight.wook.net/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

// 클라이언트가 제공한 정렬 필드가 허용 목록의 항목 중 하나와 일치하는지 확인하고,
// 일치하는 경우 정렬 필드에서 앞의 하이픈 문자(있는 경우)를 제거하여 열 이름을 추출합니다.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func ValidateFilter(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "0 보다 커야만 합니다")
	v.Check(f.Page <= 10_000_000, "page", "최대 1000만 까지 요청할 수 있습니다")
	v.Check(f.PageSize > 0, "page_size", "0 보다 커야만 합니다")
	v.Check(f.PageSize <= 100, "page_size", "최대 100 까지 요청할 수 있습니다")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "잘못된 sort 값 입니다")
}
