package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.wook.net/internal/data"
)

const version = "1.0.0"

// 연결 풀에 대한 구성 설정을 저장하기 위해 maxOpenConns,
// maxIdleConns 및 maxIdleTime 필드를 추가합니다.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production")
	// db-dsn 명령줄 플래그의 DSN 값을 config 구조체로 읽습니다.
	// 플래그가 제공되지 않으면 기본적으로 development DSN을 사용합니다.
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// 명령줄 플래그에서 연결 풀 설정을 config 구조체로 읽습니다.
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// config 구조체를 전달하면서 openDB() 도우미 함수(아래 참조)를 호출하여 연결 풀을 생성합니다.
	// 오류가 반환되면 오류를 기록하고 즉시 애플리케이션을 종료합니다.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("데이터베이스 연결 풀 설정됨")

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

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
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
