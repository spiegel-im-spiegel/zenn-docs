//go:build run
// +build run

package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		log.Println(err)
		return exitcode.Abnormal
	}
	defer db.Close()

	//Query
	_, err = db.Query("SELECT * FROM tablename") // "tablename" is not exist
	if err != nil {
		log.Println(err) // Output: pq: relation "tablename" does not exist
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
