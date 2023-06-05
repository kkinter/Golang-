package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Database struct {
	Client *sql.DB
}

func NewDatabase() (*Database, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	dbConn, err := sql.Open("pgx", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to the database: %w", err)
	}

	return &Database{Client: dbConn}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.PingContext(ctx)
}
