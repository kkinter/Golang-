package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.wook.net/internal/data"
	"greenlight.wook.net/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	// 초당 요청 수와 버스트 값에 대한 필드와 속도 제한을 활성화/비활성화하는
	// 데 사용할 수 있는 bool 필드를 포함하는 새로운 리미터 구조체를 추가합니다.
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// logger  필드를 *log.Logger 대신 *jsonlog.Logger 유형으로 변경합니다.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// 설정 값을 구성 구조체로 읽어들이는 명령줄 플래그를 생성합니다.
	// 'enabled' 설정의 기본값으로 참을 사용하는 것이 보이시죠?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "limiter 초당 최대 요청 수")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "limiter 최대 버스트")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "limiter 활성화")

	flag.Parse()

	// INFO 심각도 수준 이상의 모든 메시지를 표준 출력 스트림에 기록하는 새 jsonlog.Logger를 초기화합니다.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		// PrintFatal() 메서드를 사용하여 치명적 수준에서 오류가 포함된 로그 항목을 작성하고 종료합니다.
		// 로그 항목에 포함할 추가 프로퍼티가 없으므로 두 번째 매개변수로 nil을 전달합니다.
		logger.PrintFatal(err, nil)
	}

	// 마찬가지로 PrintInfo() 메서드를 사용하여 INFO 레벨에 메시지를 작성합니다.
	logger.PrintInfo("데이터베이스 연결 풀 설정됨", nil)

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 다시 PrintInfo() 메서드를 사용하여 INFO 수준에서 "서버 시작" 메시지를 작성합니다.
	// 하지만 이번에는 추가 속성(운영 환경 및 서버 주소)이 포함된 맵을 마지막 매개변수로 전달합니다.
	logger.PrintInfo("서버 시작", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	// PrintFatal() 메서드를 사용하여 오류를 기록하고 종료합니다.
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// 풀에서 열려 있는(사용 중 + 유휴) 연결의 최대 수를 설정합니다.
	//  0보다 작은 값을 전달하면 제한이 없음을 의미
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// 풀의 최대 유휴 연결 수를 설정합니다.
	// 0보다 작은 값을 전달하면 제한이 없다는 의미
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// time.ParseDuration() 함수를 사용하여 idle timeout duration
	// 문자열을 time.Duration 유형으로 변환합니다.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	// maximum idle timeout 설정.
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
