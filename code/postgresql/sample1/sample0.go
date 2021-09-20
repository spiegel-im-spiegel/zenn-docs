//go:build run
// +build run

package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
	logger := zerolog.New(
		os.Stdout,
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	logger.Err(os.ErrInvalid).Send()
	logger.Error().Interface("error", errs.Wrap(os.ErrInvalid)).Msg("sample error")
	return exitcode.Normal
}

func main() {
	Run().Exit()
}
