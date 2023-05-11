package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("runtime format이 잘못되었습니다")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// Runtime  유형에 UnmarshalJSON() 메서드를 구현하여 json.Unmarshaler
// 인터페이스를 만족하도록 합니다.
// 중요: UnmarshalJSON()은 receiver (Runtime 유형)를 수정해야 하므로,
// 제대로 작동하려면 포인터 receiver 를 사용해야 합니다.
// 그렇지 않으면 복사본만 수정하게 됩니다(이 메서드가 반환될 때 버려집니다).
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// 들어오는 JSON 값이 "<runtime> mins" 형식의 문자열이 될 것으로 예상하고,
	// 가장 먼저 해야 할 일은 이 문자열에서 주변 큰따옴표를 제거하는 것입니다.
	// 따옴표를 제거할 수 없으면 ErrInvalidRuntimeFormat 오류가 반환됩니다.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// 문자열을 분할하여 숫자가 포함된 부분을 분리합니다.
	parts := strings.Split(unquotedJSONValue, " ")

	// 문자열의 각 부분을 검사하여 예상되는 형식인지 확인합니다.
	// 그렇지 않은 경우 ErrInvalidRuntimeFormat 오류를 다시 반환합니다.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	// 그렇지 않으면 숫자가 포함된 문자열을 int32로 구문 분석합니다.
	// 다시 실패하면 ErrInvalidRuntimeFormat 오류를 반환합니다.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	// int32를 Runtime 타입으로 변환하고 이를 수신자에게 할당합니다.
	// 포인터의 기본값을 설정하기 위해 * 연산자를 사용하여 receiver
	// (Runtime 타입에 대한 포인터)를 참조한다는 점에 유의하세요.
	*r = Runtime(i)

	return nil
}
