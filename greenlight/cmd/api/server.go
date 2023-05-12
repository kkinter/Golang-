package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// shutdownError 채널을 생성합니다. 이를 사용하여 graceful Shutdown()
	// 함수가 반환하는 모든 오류를 수신할 것입니다.
	shutdownError := make(chan error)

	go func() {
		// 전과 같이 signals 을 가로 챕니다.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		// // 로그 항목을 "caught signal" 대신 "shutting down server"로 업데이트합니다.
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})
		// 20초 타임아웃으로 컨텍스트를 생성합니다.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		// 방금 만든 컨텍스트를 전달하면서 서버에서 Shutdown()을 호출합니다.
		// Shutdown()은 정상 종료에 성공했거나 오류(리스너를 닫는 데 문제가 있거나
		// 20초의 컨텍스트 기한이 지나기 전에 종료가 완료되지 않아서 발생할 수 있음)가
		// 발생하면 nil을 반환합니다. 이 반환값을 shutdownError 채널에 전달합니다.
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// 서버에서 Shutdown()을 호출하면 ListenAndServe()가 즉시 http.ErrServerClosed 오
	// 류를 반환합니다. 따라서 이 오류가 표시되면 실제로는 좋은 일이며 정상 종료가 시작되었다는 표시입니다.
	// 따라서 이 오류를 구체적으로 확인하여 http.ErrServerClosed가 아닌 경우에만 오류를 반환합니다.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// 그렇지 않으면 shutdownError 채널의 Shutdown()에서 반환값을 받을 때까지 기다립니다.
	// 반환값이 에러인 경우, 정상 종료에 문제가 있다는 것을 알고 에러를 반환합니다.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// 이 시점에서 정상 종료가 성공적으로 완료되었음을 알 수 있으며 "stopped server" 메시지가 기록됩니다.
	app.logger.PrintInfo("stoopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
