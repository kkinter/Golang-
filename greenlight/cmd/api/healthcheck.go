package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	// 응답에 대한 데이터가 포함된 envelope 맵을 선언합니다.
	// 이렇게 구성한 방식은 이제 환경 및 버전 데이터가 JSON 응답의
	// system_info 키 아래에 중첩된다는 것을 의미합니다.
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "서버에 문제가 발생하여 요청을 처리할 수 없습니다.", http.StatusInternalServerError)
	}
}
