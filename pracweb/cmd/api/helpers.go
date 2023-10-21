package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// write json helpers
func (app *application) writeJSON(w http.ResponseWriter, statusCode int, data any, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("요청 본문에 잘못된 형식의 JSON이 포함되어 있습니다(%d 번째 문자에서)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("요청 본문에 잘못된 형식의 JSON이 포함되어 있습니다")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("요청 본문의 %q 필드에 대한 잘못된 JSON 유형이 있습니다", unmarshalTypeError.Field)
			}
			return fmt.Errorf("요청 본문에 잘못된 JSON 형식이 있습니다(%d 번째 문자에서)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("요청 본문은 비어있지 않아야 합니다")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("요청 본문에 알 수 없는 키 %s가 포함되어 있습니다", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("본문은 %d bytes보다 크지 않아야 합니다", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("요청 본문은 단일 JSON 값만 포함해야 합니다")
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, r *http.Request, statusCode int, message any) {
	type jsonError struct {
		Message any `json:"message"`
	}

	msg := jsonError{
		Message: message,
	}

	err := app.writeJSON(w, statusCode, msg, nil)
	if err != nil {
		app.logger.Error("write json error", zap.Error(err))
		w.WriteHeader(500)
	}
}
