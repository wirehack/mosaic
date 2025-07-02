package c

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DB() *pgxpool.Pool {

	db := os.Getenv("DB")

	if db == "" {
		panic("db environment variable is not set")
	}

	pool, err := pgxpool.New(context.Background(), db)

	if err != nil {
		panic(err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	return pool
}
