package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

// New는 empty errors map.으로 새 유효성 검사기 인스턴스를 생성하는 헬퍼입니다.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError는 맵에 오류 메시지를 추가합니다
// (주어진 키에 대한 항목이 이미 존재하지 않는 한).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// 유효성 검사가 '확인'이 아닌 경우에만 맵에 오류 메시지를 추가합니다.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// 특정 값이 목록에 있으면 참을 반환하는 Generic 함수입니다.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// 문자열 값이 특정 정규식 패턴과 일치하면 참을 반환합니다.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// 슬라이스의 모든 값이 고유한 경우 참을 반환하는 Generic  함수입니다.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
