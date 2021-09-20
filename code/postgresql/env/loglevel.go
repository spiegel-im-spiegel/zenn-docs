package env

import (
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

type LoggerLevel int

const (
	LevelNop LoggerLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var levelMap = map[LoggerLevel]string{
	LevelNop:   "nop",
	LevelError: "error",
	LevelWarn:  "warn",
	LevelInfo:  "info",
	LevelDebug: "debug",
}

var zerologLevelMap = map[LoggerLevel]zerolog.Level{
	LevelNop:   zerolog.NoLevel,
	LevelError: zerolog.ErrorLevel,
	LevelWarn:  zerolog.WarnLevel,
	LevelInfo:  zerolog.InfoLevel,
	LevelDebug: zerolog.DebugLevel,
}

var pgxlogLevelMap = map[LoggerLevel]pgx.LogLevel{
	LevelNop:   pgx.LogLevelNone,
	LevelError: pgx.LogLevelError,
	LevelWarn:  pgx.LogLevelWarn,
	LevelInfo:  pgx.LogLevelInfo,
	LevelDebug: pgx.LogLevelDebug,
}

func getLogLevel(s string) LoggerLevel {
	for k, v := range levelMap {
		if strings.EqualFold(v, s) {
			return k
		}
	}
	return LevelInfo
}

func (lvl LoggerLevel) ZerlogLevel() zerolog.Level {
	if l, ok := zerologLevelMap[lvl]; ok {
		return l
	}
	return zerolog.InfoLevel
}

func (lvl LoggerLevel) PgxLogLevel() pgx.LogLevel {
	if l, ok := pgxlogLevelMap[lvl]; ok {
		return l
	}
	return pgx.LogLevelInfo
}
