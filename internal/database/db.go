package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(connStr string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}

	// test connection
	if err := Pool.Ping(context.Background()); err != nil {
		return err
	}

	log.Println("connected to postgres")
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
