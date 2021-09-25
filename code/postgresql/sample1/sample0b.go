//go:build run
// +build run

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
	file, err := os.Create("access.log")
	if err != nil {
		fmt.Println(err)
		return exitcode.Abnormal
	}
	logger := zerolog.New(
		io.MultiWriter(file, zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}),
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	logger.Err(os.ErrInvalid).Send()
	logger.Error().Interface("error", errs.Wrap(os.ErrInvalid)).Msg("sample error")
	return exitcode.Normal
}

func main() {
	Run().Exit()
}
