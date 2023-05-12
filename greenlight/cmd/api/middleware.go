package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer 함수를 생성합니다(이 함수는 패닉이 발생했을 때 Go가 스택을 풀 때 항상 실행됩니다).
		defer func() {
			// builtin recover function을 사용하여 패닉이 발생했는지 여부를 확인합니다.
			if err := recover(); err != nil {
				// 패닉이 발생한 경우 응답에 "Connection: close"" 헤더를 설정합니다.
				// 이는 응답이 전송된 후 Go의 HTTP 서버가 현재 연결을 자동으로 닫도록 하는 트리거 역할을 합니다.
				w.Header().Set("Connection", "close")
				// recover()가 반환하는 값의 유형이 any이므로 fmt.Errorf()를 사용하여 오류로 정규화하고
				// serverErrorResponse() 헬퍼를 호출합니다. 그러면 오류 수준에서
				//사용자 정의 로거 유형을 사용하여 오류를 기록하고 클라이언트에게 500 내부 서버 오류 응답을 보냅니다.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// 초당 평균 2건의 요청을 허용하고 단일 '버스트'에 최대 4건의
	// 요청을 허용하는 새로운 전송률 제한기를 초기화합니다.

	limiter := rate.NewLimiter(2, 4)

	// 우리가 반환하는 함수는 리미터 변수를 'closes over' 클로저입니다.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// limiter.Allow()을 호출하여 요청이 허용되는지 확인하고,
		// 허용되지 않는 경우 rateLimitExceededResponse() 헬퍼를 호출하여
		// 429 너무 많은 요청 응답을 반환합니다(이 헬퍼는 잠시 후에 생성할 예정입니다).
		if !limiter.Allow() {
			app.rateLimitExccededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
