package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// 애플리케이션 버전 번호가 포함된 문자열을 선언합니다.
// 이 책의 뒷부분에서는 빌드 시점에 자동으로 생성하지만
// 지금은 버전 번호를 하드코딩된 전역 상수로 저장하겠습니다.
const version = "1.0.0"

// 애플리케이션의 모든 구성 설정을 저장할 구성 구조체를 정의합니다.
// 현재로서는 서버가 수신 대기할 네트워크 포트와 애플리케이션의 현재 운영 환경 이름(개발 및 스테이징)만 구성 설정에 포함됩니다.
// 애플리케이션의 현재 운영 환경 이름(개발, 스테이징, 프로덕션 등)입니다.
// 애플리케이션이 시작될 때 명령줄 플래그에서 이러한 구성 설정을 읽어들일 것입니다.
type config struct {
	port int
	env  string
}

// HTTP 핸들러, 헬퍼, 미들웨어에 대한 종속성을 담을 애플리케이션 구조체를 정의합니다.
// 현재는 구성 구조체와 로거의 복사본만 포함되어 있지만
// 빌드가 진행됨에 따라 더 많은 것을 포함하도록 확장할 것입니다.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// config 구조체의 인스턴스를 선언합니다.
	var cfg config

	// port  및 env 명령줄 플래그 값을 config 구조체에 읽어들입니다.
	//해당 플래그가 제공되지 않으면 기본적으로 포트 번호 4000과
	// environment "development"을 사용합니다.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production")
	flag.Parse()

	// 현재 날짜와 시간이 접두사로 붙은 메시지를 표준 출력 스트림에 기록하는
	// 새 로거를 초기화합니다.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// config 구조체와 logger 포함하는 애플리케이션 구조체의 인스턴스를 선언합니다.
	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
