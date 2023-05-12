package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"greenlight.wook.net/internal/data"
	"greenlight.wook.net/internal/jsonlog"
	"greenlight.wook.net/internal/mailer"
)

const version = "1.0.0"

// SMTP 서버 설정을 보유하도록 config 구조를 업데이트합니다.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// 애플리케이션 구조체에 sync.WaitGroup을 포함합니다. sync.WaitGroup 유형의 0 값은 '카운터' 값이 0인 유효하고
//
//	사용 가능한 sync.WaitGroup이므로 사용하기 전에 초기화하기 위해 다른 작업을 할 필요가 없습니다.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "limiter 초당 최대 요청 수")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "limiter 최대 버스트")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "limiter 활성화")

	// 메일트랩 설정을 기본값으로 사용하여 SMTP 서버 구성 설정을 구성 구조체에 읽어들입니다.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "b66e54365d5d63", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "ba2b076e007424", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.wook.net>", "SMTP sender")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("데이터베이스 연결 풀 설정됨", nil)

	// 명령줄 플래그의 설정을 사용하여 새 메일러 인스턴스를 초기화하고 애플리케이션 구조에 추가합니다.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
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
