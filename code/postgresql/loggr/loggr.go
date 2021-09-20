package loggr

import (
	"fmt"
	"io"
	"os"
	"sample/env"
	"time"

	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/cache"
)

func New() *zerolog.Logger {
	logger := zerolog.Nop()
	if env.ZerologLevel() == zerolog.NoLevel {
		return &logger
	}

	// enable ConsoleWriter
	var stdout io.Writer = os.Stdout
	if env.ZerologTerminal() {
		stdout = zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}
	}

	// make path to ${XDG_CACHE_HOME}/${ServiceName}/access.YYYYMMDD.log file and create logger
	logpath := cache.Path(env.ServiceName, fmt.Sprintf("access.%s.log", time.Now().Local().Format("20060102")))
	file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		logger = zerolog.New(stdout)
	} else {
		logger = zerolog.New(io.MultiWriter(
			file,
			stdout,
		))
	}
	logger = logger.Level(env.ZerologLevel()).With().Timestamp().Logger()

	if err != nil {
		logger.Error().Interface("error", errs.Wrap(err, errs.WithContext("logpath", logpath))).Str("logpath", logpath).Msg("error in opening logfile")
	}
	return &logger
}
