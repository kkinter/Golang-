package data

import (
	"fmt"
	"strconv"
)

// 기본 유형이 int32(Movie 구조체 필드와 동일)인 사용자 정의 Runtime 유형을 선언합니다.
type Runtime int32

// 런타임 유형에 json.Marshaler 인터페이스를 만족하도록 MarshalJSON() 메서드를 구현합니다.
// 이 메서드는 동영상 런타임에 대한 JSON 인코딩된 값을 반환해야 합니다(이 경우 "<런타임> 분" 형식의 문자열을 반환합니다).
func (r Runtime) MarshalJSON() ([]byte, error) {
	// 필요한 형식으로 동영상 런타임이 포함된 문자열을 생성합니다.
	jsonValue := fmt.Sprintf("%d mins", r)
	// 문자열을 큰따옴표로 묶으려면 문자열에 strconv.Quote() 함수를 사용합니다. 유효한 *JSON 문자열*이 되려면 큰따옴표로 묶어야 합니다.
	quotedJSONValue := strconv.Quote(jsonValue)

	// 따옴표로 묶은 문자열 값을 바이트 슬라이스로 변환하여 반환합니다.
	return []byte(quotedJSONValue), nil
}
