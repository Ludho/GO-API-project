package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Url string
}

func connectToDB(dbConfig *DBConfig) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbConfig.Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error when parsing the db URL connection string : %v", err)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		fmt.Printf("connection to db successfull")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Errorf("connection to db failed :%v", err))
	}

	return pool
}
