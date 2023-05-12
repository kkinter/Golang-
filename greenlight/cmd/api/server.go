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
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// 이전처럼 서버에서 Shutdown()을 호출하지만, 이제 오류를 반환하는 경우에만
		// shutdownError 채널에 전송합니다.
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// 백그라운드 고루틴이 작업을 완료하기를 기다리고 있다는 메시지를 기록합니다.
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		// Wait()를 호출하여 WaitGroup 카운터가 0이 될 때까지 차단하고, 백그라운드 고루틴이
		// 완료될 때까지 차단합니다. 그런 다음 종료가 문제 없이 완료되었음을 나타내기 위해
		// shutdownError 채널에 nil을 반환합니다.
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stoopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
