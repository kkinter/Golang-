package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"greenlight.wook.net/internal/validator"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// 요청 본문의 크기를 1MB로 제한하려면 http.MaxBytesReader()를 사용합니다.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// json.Decoder를 초기화하고, 디코딩하기 전에 DisallowUnknownFields()
	// 메서드를 호출합니다. 즉, 이제 클라이언트의 JSON에 대상 대상에 매핑할 수 없는
	// 필드가 포함된 경우 디코더가 해당 필드를 무시하는 대신 오류를 반환합니다.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// 요청 본문을 대상으로 디코딩합니다.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarhalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		// maxBytesError  변수를 새로 추가합니다.
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("본문에 잘못된 형식의 JSON이 포함되어 있습니다(%d 문자에서)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("본문에 잘못된 형식의 JSON이 포함되어 있습니다")

		case errors.As(err, &unmarhalTypeError):
			if unmarhalTypeError.Field != "" {
				return fmt.Errorf("본문에 %q 필드에 대한 잘못된 JSON 유형이 있습니다", unmarhalTypeError.Field)
			}
			return fmt.Errorf("본문에 잘못된 JSON 형식이 있습니다(%d 문자에서)", unmarhalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("본문이 비어 있지 않아야 합니다")

		// JSON에 대상 대상에 매핑할 수 없는 필드가 포함된 경우 이제 Decode()는
		// "json: 알 수 없는 필드 "<이름>"" 형식의 오류 메시지를 반환합니다.
		// 이를 확인하고 오류에서 필드 이름을 추출하여 사용자 정의 오류 메시지에
		// 보간합니다. 향후 이를 별도의 오류 유형으로 전환하는 것과 관련하여
		// https://github.com/golang/go/issues/29035 에 미해결 이슈가 있다는 점에 유의하세요.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("본문에 알 수 없는 키 %s가 포함되어 있습니다", fieldName)
		// errors.As() 함수를 사용하여 오류의 유형이 *http.MaxBytesError인지 확인합니다.
		// 만약 그렇다면 요청 본문이 크기 제한인 1MB를 초과했음을 의미하며
		// 명확한 오류 메시지를 반환합니다.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("본문은 %d bytes보다 크지 않아야 합니다", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// 빈 익명 구조체에 대한 포인터를 대상으로 사용하여 Decode()를 다시 호출합니다.
	// 요청 본문에 단일 JSON 값만 포함된 경우 io.EOF 오류가 반환됩니다.
	// 따라서 다른 값을 받으면 요청 본문에 추가 데이터가 있다는 것을 알고
	// 자체 사용자 정의 오류 메시지를 반환합니다.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("본문은 단일 JSON 값만 포함해야 합니다")
	}

	return nil
}

// readString() 헬퍼는 쿼리 문자열에서 문자열 값을 반환하거나
// 일치하는 키를 찾을 수 없는 경우 제공된 기본값을 반환합니다.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {

	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// readCSV() 헬퍼는 쿼리 문자열에서 문자열 값을 읽은 다음 쉼표
// 문자의 조각으로 분할합니다. 일치하는 키를 찾을 수 없으면 제공된 기본값을 반환합니다.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// 쿼리 문자열에서 값을 추출합니다.
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// readInt() 헬퍼는 쿼리 문자열에서 문자열 값을 읽고 정수로 변환한 후 반환합니다.
// 일치하는 키를 찾을 수 없으면 제공된 기본값을 반환합니다.
// 값을 정수로 변환할 수 없는 경우 제공된 유효성 검사기 인스턴스에 오류 메시지를 기록합니다.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be and integer value")
		return defaultValue
	}

	return i
}

func (app *application) background(fn func()) {
	// WaitGroup 카운터를 증가시킵니다.
	app.wg.Add(1)

	go func() {

		// defer를 사용하면 고루틴이 반환되기 전에 WaitGroup 카운터가 감소합니다.
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
