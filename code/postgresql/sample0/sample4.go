//go:build run
// +build run

package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spiegel-im-spiegel/gocli/config"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func init() {
	//load ~/.config/elephantsql/env file
	if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil {
		panic(err)
	}
}

func Run() exitcode.ExitCode {
	// connect PostgreSQL service
	conn, err := pgxpool.Connect(context.TODO(), os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		log.Println(err)
		return exitcode.Abnormal
	}
	defer conn.Close()

	_, err = conn.Query(context.TODO(), "SELECT * FROM tablename") // "tablename" is not exist
	if err != nil {
		log.Println(err) // Output: ERROR: relation "tablename" does not exist (SQLSTATE 42P01)
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
