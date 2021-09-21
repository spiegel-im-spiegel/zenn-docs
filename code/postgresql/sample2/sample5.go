//go:build run
// +build run

package main

import (
	"fmt"
	"os"
	"sample/files"
	"sample/gorm/model"
	"sample/orm"

	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
	"gorm.io/gorm"
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	file1 := "files/file1.txt"
	bin1, err := files.GetBinary(file1)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	file2 := "files/file2.txt"
	bin2, err := files.GetBinary(file2)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}

	data := &model.User{
		Username: "Alice",
		BinaryFiles: []model.BinaryFile{
			{Filename: file1, Body: bin1},
			{Filename: file2, Body: bin2},
		},
	}

	// drop table
	tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Create(data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
