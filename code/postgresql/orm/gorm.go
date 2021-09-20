package orm

import (
	"sample/dbconn"
	"sample/env"

	"github.com/rs/zerolog"
	"github.com/spiegel-im-spiegel/errs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormContext struct {
	Db     *gorm.DB
	Logger *zerolog.Logger
}

func NewGORM() (*GormContext, error) {
	pgxCtx, err := dbconn.NewPgx()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	gormCtx := &GormContext{
		Logger: pgxCtx.GetLogger(),
	}
	loggr := logger.Discard
	if env.LogLevel() == env.LevelDebug {
		loggr = logger.Default
	}
	gormCtx.Db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: pgxCtx.GetDb(),
	}), &gorm.Config{
		Logger: loggr,
	})
	if err != nil {
		pgxCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Msg("error in gorm.Open() method")
		pgxCtx.Close()
		return nil, errs.Wrap(err)
	}
	return gormCtx, nil
}

func (gormCtx *GormContext) GetDb() *gorm.DB {
	if gormCtx == nil {
		return nil
	}
	return gormCtx.Db
}

func (gormCtx *GormContext) GetLogger() *zerolog.Logger {
	if gormCtx == nil {
		lggr := zerolog.Nop()
		return &lggr
	}
	return gormCtx.Logger
}

func (gormCtx *GormContext) Close() error {
	if db := gormCtx.GetDb(); db != nil {
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.Close()
		}
	}
	return nil
}
