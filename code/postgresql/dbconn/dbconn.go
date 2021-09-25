package dbconn

import (
	"database/sql"
	"sample/env"
	"sample/loggr"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
)

type PgxContext struct {
	Db     *sql.DB
	Logger *zerolog.Logger
}

func NewPgx() (*PgxContext, error) {
	dbctx := &PgxContext{
		Logger: loggr.New(),
	}
	cfg, err := pgx.ParseConfig(env.PostgresDSN())
	if err != nil {
		dbctx.Logger.Error().Interface("error", errs.Wrap(err)).Msg("error in pgx.ParseConfig() method")
		return nil, errs.Wrap(err, errs.WithContext("dsn", env.PostgresDSN()))
	}
	cfg.Logger = zerologadapter.NewLogger(*dbctx.Logger)
	cfg.LogLevel = env.PgxlogLevel()
	dbctx.Db = stdlib.OpenDB(*cfg)

	return dbctx, nil
}

func (dbctx *PgxContext) GetDb() *sql.DB {
	if dbctx == nil {
		return nil
	}
	return dbctx.Db
}

func (dbctx *PgxContext) GetLogger() *zerolog.Logger {
	if dbctx == nil {
		lggr := zerolog.Nop()
		return &lggr
	}
	return dbctx.Logger
}

func (dbctx *PgxContext) Acquire() (*pgx.Conn, error) {
	if db := dbctx.GetDb(); db != nil {
		conn, err := stdlib.AcquireConn(db)
		return conn, errs.Wrap(err)
	}
	return nil, errs.New("*sql.DB instance is nil.")
}

func (dbctx *PgxContext) Close() error {
	if db := dbctx.GetDb(); db != nil {
		return errs.Wrap(db.Close())
	}
	return nil
}
