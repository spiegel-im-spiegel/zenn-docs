//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/cache"
	"github.com/spiegel-im-spiegel/gocli/config"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func init() {
	//load ~/.config/elephantsql/env file
	if err := godotenv.Load(config.Path("elephantsql", "env")); err != nil {
		panic(err)
	}
}

func CreateLogger() zerolog.Logger {
	logger := zerolog.Nop()
	logpath := cache.Path("pgx", fmt.Sprintf("access.%s.log", time.Now().Local().Format("20060102"))) // logpath = ~/.cache/pgx/access.YYYYMMDD.log
	file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		logger = zerolog.New(os.Stdout)
	} else {
		logger = zerolog.New(io.MultiWriter(
			file,
			zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
		))
	}
	logger = logger.Level(zerolog.DebugLevel).With().Timestamp().Logger()
	if err != nil {
		logger.Error().Interface("error", errs.Wrap(err, errs.WithContext("logpath", logpath))).Str("logpath", logpath).Msg("error in opening logfile")
	}
	return logger
}

func Run() exitcode.ExitCode {
	// get logger
	logger := CreateLogger()

	// connect PostgreSQL service
	cfg, err := pgx.ParseConfig(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		logger.Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	cfg.Logger = zerologadapter.NewLogger(logger)
	cfg.LogLevel = pgx.LogLevelDebug
	db := stdlib.OpenDB(*cfg)
	defer db.Close()

	// get connection from connection pool
	conn, err := db.Conn(context.TODO())
	if err != nil {
		logger.Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	defer conn.Close()

	// query
	_, err = conn.QueryContext(context.TODO(), "SELECT * FROM tablename") // "tablename" is not exist
	if err != nil {
		logger.Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
