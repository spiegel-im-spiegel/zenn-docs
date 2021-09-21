package env

import (
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

var levelMap = map[string]LoggerLevel{
	"nop":   LevelNop,
	"error": LevelError,
	"warn":  LevelWarn,
	"info":  LevelInfo,
	"debug": LevelDebug,
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
	if lvl, ok := levelMap[s]; ok {
		return lvl
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
