package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// 로그 항목의 심각도 수준을 나타내는 Level 유형을 정의합니다.
type Level int8

const (
	LevelInfo  Level = iota // 값이 0입니다.
	LevelError              // 값이 1입니다.
	LevelFatal              // 값이 2입니다.
	LevelOff                // 값이 3입니다.
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// 사용자 정의 로거 유형을 정의합니다. 여기에는 로그 항목이
// 기록될 출력 대상, 로그 항목이 기록될 최소 심각도 수준,
// 쓰기 조정을 위한 뮤텍스가 저장됩니다.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// 최소 심각도 수준 이상의 로그 항목을 특정 출력 대상에 기록하는 새 Logger 인스턴스를 반환합니다.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// 다양한 수준에서 로그 항목을 작성하기 위한 몇 가지 헬퍼 메서드를 선언합니다.
// 이 메서드들은 모두 로그 항목에 표시하려는 임의의 '속성'을 포함할 수 있는
//
//	맵을 두 번째 매개변수로 받습니다.
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1) // 치명적인 수준의 항목에 대해서는 애플리케이션도 종료합니다.
}

// Print는 로그 항목을 작성하는 내부 메서드입니다.
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}
	//  로그 항목의 데이터를 저장하는 익명 구조체를 선언합니다.
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	// 오류 및 치명적 수준의 항목에 대한 스택 추적을 포함합니다.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	// 익명 구조체를 JSON으로 마샬링하여 line 변수에 저장합니다.
	// JSON을 생성하는 데 문제가 있는 경우 로그 항목의 내용을 일반 텍스트 오류 메시지로 대신 설정하세요.
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}
	// 출력 대상에 대한 두 개의 쓰기가 동시에 일어나지 않도록 뮤텍스를 잠급니다.
	// 이렇게 하지 않으면 두 개 이상의 로그 항목에 대한 텍스트가 출력에 섞일 수 있습니다.
	l.mu.Lock()
	defer l.mu.Unlock()

	// 로그 항목 뒤에 줄 바꿈으로 로그 항목을 작성합니다.
	return l.out.Write(append(line, '\n'))
}

// 또한 Logger 유형에 Write() 메서드를 구현하여 io.Writer 인터페이스를 만족하도록 합니다.
// 이렇게 하면 추가 프로퍼티 없이 오류 수준에서 로그 항목을 기록합니다.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
