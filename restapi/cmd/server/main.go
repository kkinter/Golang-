package main

import (
	"context"
	"fmt"

	"github.com/go-api-prac/internal/comment"
	"github.com/go-api-prac/internal/db"
)

// Run - go 애플리케이션의 인스턴스화 및 시작을 담당합니다.
func Run() error {
	fmt.Println("starting up our app")

	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the db")
		return err
	}
	// if err := db.Ping(context.Background()); err != nil {
	// 	return err
	// }
	if err := db.MigrateDB(); err != nil {
		fmt.Println("failed to migrate db")
		return err
	}

	fmt.Println("successfully connected and pinged db")

	cmtService := comment.NewService(db)
	fmt.Println(cmtService.GetComment(
		context.Background(),
		"206c6d84-0389-11ee-be56-0242ac120002",
	))

	return nil
}

func main() {
	fmt.Println("GO REST API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
