package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.wook.net/internal/data"
	"greenlight.wook.net/internal/validator"
)

func (app *application) creteMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResopnse(w, r, err)
		return
	}

	// movie 변수에는 Movie 구조체에 대한 *포인터*가 포함되어 있습니다.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// 새 유효성 검사기를 초기화합니다.
	v := validator.New()

	// 검사 중 하나라도 실패하면 ValidateMovie() 함수를 호출하고
	// 오류가 포함된 응답을 반환합니다.
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 유효성이 검사된 movies 구조체에 대한 포인터를 전달하여 movies 모델에서
	// Insert() 메서드를 호출합니다. 그러면 데이터베이스에 레코드가 생성되고
	// 시스템에서 생성된 정보로 movie 구조체가 업데이트됩니다.
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// HTTP 응답을 보낼 때 클라이언트가 새로 생성된 리소스를 찾을 수 있는 URL을 알 수 있도록,
	// Location 헤더를 포함하려고 합니다. 빈 http.Header 맵을 만든 다음
	// Set() 메서드를 사용하여 새 Location 헤더를 추가하고 URL에 새 movie 에
	// 대해 시스템에서 생성된 ID를 보간합니다.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// 201 Created 상태 코드, 응답 본문의 movie 데이터, Location 헤더가
	// 포함된 JSON 응답을 작성합니다.
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Get() 메서드를 호출하여 특정 동영상에 대한 데이터를 가져옵니다.
	// 또한 errors.Is() 함수를 사용하여 data.ErrRecordNotFound 오류를 반환하는지
	// 확인해야 하며, 이 경우 클라이언트에 404 찾을 수 없음 응답을 전송합니다.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 제목, 연도 및 런타임 필드에 포인터를 사용합니다.
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	// var input struct {
	// 	Title   string       `json:"title"`
	// 	Year    int32        `json:"year"`
	// 	Runtime data.Runtime `json:"runtime"`
	// 	Genres  []string     `json:"genres"`
	// }

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResopnse(w, r, err)
		return
	}

	// input.Title 값이 nil이면 JSON 요청 본문에서 해당 "title" 키/값 쌍이 제공되지 않았다는
	// 것을 알 수 있습니다. 따라서 동영상 레코드를 변경하지 않고 그대로 둡니다.
	// 그렇지 않으면 새 제목 값으로 동영상 레코드를 업데이트합니다.
	// 중요한 점은 input.Title이 이제 문자열에 대한 포인터이기 때문에 동영상 레코드에
	// 할당하기 전에 * 연산자를 사용하여 포인터를 역참조하여 기본 값을 가져와야 한다는 것입니다.
	if input.Title != nil {
		movie.Title = *input.Title
	}

	// input 구조체의 다른 필드에 대해서도 동일한 작업을 수행합니다.
	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 모든 ErrEditConflict 오류를 가로채고 새로운 editConflictResponse() 헬퍼를 호출합니다.
	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilter(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 메타데이터 구조체를 반환값으로 받습니다.
	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 응답 envelope에 메타데이터를 포함합니다.
	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
