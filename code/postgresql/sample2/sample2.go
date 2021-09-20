//go:build run
// +build run

package main

import (
	"fmt"
	"os"
	"sample/orm"

	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	// query
	var results []map[string]interface{}
	tx := gormCtx.GetDb().Table("tablename").Find(&results) // "tablename" is not exist
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
