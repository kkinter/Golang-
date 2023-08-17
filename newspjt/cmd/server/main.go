package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	os.Exit(start())
}

func start() int {
	logEnv := getStringOrDefault("LOG_ENV", "development")
	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("로거 설정에 오류가 발생했습니다:", err)
		return 1
	}

	defer func() {
		_ = log.Sync()
	}()

	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func getStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}
