package data

import (
	"math"
	"strings"

	"greenlight.wook.net/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calculateMetadata() 함수는 총 레코드 수, 현재 페이지 및 페이지 크기 값이
// 주어지면 적절한 페이지 매김 메타데이터 값을 계산합니다. 마지막 페이지 값은
// 부동 소수점을 가장 가까운 정수로 반올림하는 math.Ceil() 함수를 사용하여
// 계산된다는 점에 유의하세요. 예를 들어 총 레코드 수가 12개이고 페이지 크기가 5인 경우
// 마지막 페이지 값은 math.Ceil(12/5) = 3이 됩니다.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	// 레코드가 없는 경우 빈 메타데이터 구조체를 반환한다는 점에 유의하세요.
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
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

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilter(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "0 보다 커야만 합니다")
	v.Check(f.Page <= 10_000_000, "page", "최대 1000만 까지 요청할 수 있습니다")
	v.Check(f.PageSize > 0, "page_size", "0 보다 커야만 합니다")
	v.Check(f.PageSize <= 100, "page_size", "최대 100 까지 요청할 수 있습니다")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "잘못된 sort 값 입니다")
}
