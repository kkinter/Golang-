package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := {"error": message}
	// writeJSON() 헬퍼를 사용하여 응답을 작성합니다. 이 과정에서 오류가 반환되면
	// 이를 기록하고 클라이언트에 500 내부 서버 오류 상태 코드가 포함된
	// 빈 응답을 보내는 것으로 되돌아갑니다.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}

}
