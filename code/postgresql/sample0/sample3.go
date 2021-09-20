//go:build run
// +build run

package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/spiegel-im-spiegel/gocli/config"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
	// get PostgreSQL connection URL
	if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil { //load ~/.config/elephantsql/env file
		log.Println(err)
		return exitcode.Abnormal
	}

	// connect PostgreSQL service
	conn, err := pgx.Connect(context.TODO(), os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		log.Println(err)
		return exitcode.Abnormal
	}
	defer conn.Close(context.TODO())

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
