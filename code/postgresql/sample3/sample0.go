//go:build run
// +build run

package main

import (
	"os"
	"sample/loggr"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/examples/start/ent"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
	// get logger
	zlogger := loggr.New()

	// create ent.Client instance for PostgreSQL service
	cfg, err := pgx.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		zlogger.Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	cfg.Logger = zerologadapter.NewLogger(*zlogger)
	cfg.LogLevel = pgx.LogLevelDebug
	client := ent.NewClient(
		ent.Driver(
			sql.OpenDB(dialect.Postgres, stdlib.OpenDB(*cfg)),
		),
	)
	defer client.Close()

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
