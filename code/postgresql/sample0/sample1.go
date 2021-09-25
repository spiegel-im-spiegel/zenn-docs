//go:build run
// +build run

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spiegel-im-spiegel/gocli/config"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
	// get PostgreSQL connection URL
	if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil { // load ~/.config/elephantsql/env file
		log.Println(err)
		return exitcode.Abnormal
	}
	fmt.Println(os.Getenv("ELEPHANTSQL_URL")) // Output: postgres://username:password@hostname:port/databasename
	return exitcode.Normal
}

func main() {
	Run().Exit()
}
