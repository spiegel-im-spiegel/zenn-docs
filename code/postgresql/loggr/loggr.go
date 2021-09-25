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
	if env.EnableLogFile() {
		// make path to ${XDG_CACHE_HOME}/${ServiceName}/access.YYYYMMDD.log file and create logger
		logpath := cache.Path(env.ServiceName, fmt.Sprintf("access.%s.log", time.Now().Local().Format("20060102")))
		if file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
			logger = zerolog.New(os.Stdout).Level(env.ZerologLevel()).With().Timestamp().Logger()
			logger.Error().Interface("error", errs.Wrap(err)).Str("logpath", logpath).Msg("error in opening logfile")
		} else {
			logger = zerolog.New(io.MultiWriter(
				file,
				zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
			)).Level(env.ZerologLevel()).With().Timestamp().Logger()
		}
		return &logger
	}
	logger = zerolog.New(os.Stdout).Level(env.ZerologLevel()).With().Timestamp().Logger()
	return &logger
}
