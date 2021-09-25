//go:build run
// +build run

package main

import (
	"context"
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

	// delete all data in users table
	tx := gormCtx.GetDb().WithContext(context.TODO()).Exec("DELETE FROM users") // delete permanently
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
