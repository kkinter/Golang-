package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()
	// http.HandlerFunc() 어댑터를 사용하여 notFoundResponse()
	// 헬퍼를 http.Handler로 변환한 다음 404 찾을 수 없음 응답에 대한
	// 사용자 지정 오류 핸들러로 설정합니다.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// 마찬가지로 methodNotAllowedResponse() 헬퍼를 http.Handler로 변환하고
	//  405 메서드 허용되지 않음 응답에 대한 사용자 정의 오류 처리기로 설정합니다.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.creteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	return router
}
