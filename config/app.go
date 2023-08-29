package config

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Connection *pgxpool.Pool
}

func Load() *AppConfig {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Couldn't load env vars from .env file")
	}

	// load environment VALUES to config structs
	dbConfig := DBConfig{}
	envconfig.MustProcess("DB", &dbConfig)

	app := &AppConfig{
		Connection: connectToDB(&dbConfig),
	}

	return app
}
