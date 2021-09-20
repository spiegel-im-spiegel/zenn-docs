//go:build run
// +build run

package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/cache"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

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
	logger := CreateLogger()
	logger.Err(os.ErrInvalid).Send()
	logger.Error().Interface("error", errs.Wrap(os.ErrInvalid)).Msg("sample error")
	return exitcode.Normal
}

func main() {
	Run().Exit()
}
