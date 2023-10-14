package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// jsno helper !
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("본문에 잘못된 형식의 JSON이 포함되어 있습니다(%d 문자에서)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("본문에 잘못된 형식의 JSON이 포함되어 있습니다")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("본문에 %q 필드에 대한 잘못된 JSON 유형이 있습니다", unmarshalTypeError.Field)
			}
			return fmt.Errorf("본문에 잘못된 JSON 형식이 있습니다(%d 문자에서)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("본문이 비어 있지 않아야 합니다")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}
