package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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

	// 각 클라이언트의 rate limiter와 last seen  시간을 저장할 클라이언트 구조체를 정의합니다.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu sync.Mutex
		// 값이 클라이언트 구조체에 대한 포인터가 되도록 맵을 업데이트합니다.
		clients = make(map[string]*client)
	)

	// 1분에 한 번씩 클라이언트 맵에서 오래된 항목을 제거하는 백그라운드 고루틴을 실행합니다.
	go func() {
		for {
			time.Sleep(time.Minute)
			// 정리하는 동안  rate limiter 검사가 발생하지 않도록 뮤텍스를 잠급니다.
			mu.Lock()
			// 모든 클라이언트를 반복합니다. 지난 3분 이내에 고객이 보이지 않았다면
			// map에서 해당 항목을 삭제합니다.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			// 중요한 것은 정리가 완료되면 뮤텍스의 잠금을 해제해야 한다는 것입니다.
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			// 아직 존재하지 않는 경우 새 클라이언트 구조체를 생성하고 맵에 추가합니다.
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// 클라이언트의 마지막으로 본 시간을 업데이트합니다.
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExccededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
